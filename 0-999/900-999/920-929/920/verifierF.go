package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func divCount(x int) int {
	cnt := 0
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			cnt += 2
			if d*d == x {
				cnt--
			}
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(6)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(20) + 1
		}
		queries := make([][3]int, m)
		sumQueries := 0
		for i := 0; i < m; i++ {
			typ := rand.Intn(2) + 1
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			queries[i] = [3]int{typ, l, r}
			if typ == 2 {
				sumQueries++
			}
		}
		if sumQueries == 0 {
			queries[0][0] = 2
		}
		// compute expected answers
		expected := []int{}
		arrCopy := append([]int(nil), arr...)
		for _, q := range queries {
			if q[0] == 1 {
				for i := q[1] - 1; i < q[2]; i++ {
					arrCopy[i] = divCount(arrCopy[i])
				}
			} else {
				total := 0
				for i := q[1] - 1; i < q[2]; i++ {
					total += arrCopy[i]
				}
				expected = append(expected, total)
			}
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for i, v := range arr {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
		for _, q := range queries {
			fmt.Fprintf(&input, "%d %d %d\n", q[0], q[1], q[2])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to run binary:", err)
			os.Exit(1)
		}
		parts := strings.Fields(string(out))
		if len(parts) != len(expected) {
			fmt.Printf("Test %d failed: expected %d outputs got %d\n", t+1, len(expected), len(parts))
			os.Exit(1)
		}
		for i, exp := range expected {
			val, err := strconv.Atoi(parts[i])
			if err != nil || val != exp {
				fmt.Printf("Test %d output %d mismatch: expected %d got %s\n", t+1, i+1, exp, parts[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

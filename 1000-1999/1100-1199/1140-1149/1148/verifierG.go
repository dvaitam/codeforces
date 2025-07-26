package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase() (string, []int) {
	k := rand.Intn(3) + 2
	n := rand.Intn(5) + 2*k
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(50) + 2
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), a
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isValid(set []int, vals []int) bool {
	k := len(set)
	for i := 0; i < k; i++ {
		if set[i] < 1 || set[i] > len(vals) {
			return false
		}
		for j := i + 1; j < k; j++ {
			if set[i] == set[j] {
				return false
			}
		}
	}
	allFair := true
	noneFair := true
	for i := 0; i < k; i++ {
		fair := true
		for j := 0; j < k; j++ {
			if i == j {
				continue
			}
			if gcd(vals[set[i]-1], vals[set[j]-1]) == 1 {
				fair = false
			}
		}
		if fair {
			noneFair = false
		} else {
			allFair = false
		}
	}
	return allFair || noneFair
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, vals := genCase()
		out, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		first := strings.Fields(strings.Split(strings.TrimSpace(input), "\n")[0])
		k, _ := strconv.Atoi(first[1])
		if len(fields) != k {
			fmt.Fprintf(os.Stderr, "wrong number of indices on test %d\n", t+1)
			os.Exit(1)
		}
		set := make([]int, k)
		for i := 0; i < k; i++ {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad output on test %d\n", t+1)
				os.Exit(1)
			}
			set[i] = v
		}
		if !isValid(set, vals) {
			fmt.Fprintf(os.Stderr, "invalid set on test %d\ninput:\n%soutput:\n%s", t+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

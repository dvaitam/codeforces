package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedD(n, a, b, k int, arr []int) string {
	freq := make(map[int]int, n)
	for _, r := range arr {
		freq[r]++
	}
	unique := make([]int, 0, len(freq))
	for v := range freq {
		unique = append(unique, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(unique)))

	ans := 0
	for _, r := range unique {
		if r%k != 0 {
			continue
		}
		s := r / k
		cntR := freq[r]
		cntS, ok := freq[s]
		if !ok || cntR == 0 || cntS == 0 {
			continue
		}
		t := cntR / b
		if v := cntS / a; v < t {
			t = v
		}
		if t > 0 {
			freq[r] -= t * b
			freq[s] -= t * a
			ans += t
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		a := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		k := rand.Intn(5) + 2
		arr := make([]int, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d %d\n", n, a, b, k))
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(50) + 1
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')
		expected := expectedD(n, a, b, k, append([]int(nil), arr...))
		out, err := runProgram(bin, []byte(input.String()))
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Expected:", expected)
			fmt.Println("Output:", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}

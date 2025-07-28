package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func generateTests() []int {
	r := rand.New(rand.NewSource(4))
	tests := make([]int, 100)
	for i := range tests {
		tests[i] = r.Intn(20) + 2
	}
	return tests
}

func solveOne(n int) []int64 {
	nn := int64(n)
	for sq := nn + 100; ; sq++ {
		mx := sq + 1
		ans := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			if i == n-1 {
				ans[i] = mx
			} else {
				ans[i] = int64(i + 1)
			}
			sum += ans[i]
		}
		rem := sq*sq - sum
		inc := rem / nn
		mod := rem % nn
		if mod > nn-2 {
			continue
		}
		for i := n - 1; i >= 0; i-- {
			ans[i] += inc
			if i > 0 && i+1 < n && mod > 0 {
				mod--
				ans[i]++
			}
		}
		return ans
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, n := range tests {
		fmt.Fprintln(&input, n)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	for i, n := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != n {
			fmt.Fprintf(os.Stderr, "wrong length on test %d\n", i+1)
			os.Exit(1)
		}
		want := solveOne(n)
		for j, p := range parts {
			v, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid integer on test %d\n", i+1)
				os.Exit(1)
			}
			if v != want[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All test cases passed.")
}

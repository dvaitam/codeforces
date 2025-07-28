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
	r := rand.New(rand.NewSource(2))
	tests := make([]int, 100)
	for i := range tests {
		tests[i] = r.Intn(20) + 1
	}
	return tests
}

func solve(n int) []int {
	if n == 1 {
		return []int{69}
	}
	if n == 2 {
		return []int{1, 3}
	}
	if n%2 == 1 {
		ans := make([]int, n)
		for i := range ans {
			ans[i] = 7
		}
		return ans
	}
	ans := make([]int, n)
	for i := 0; i < n-3; i++ {
		ans[i] = 2
	}
	ans[n-3] = 1
	ans[n-2] = 2
	ans[n-1] = 3
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "missing output for test case %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != n {
			fmt.Fprintf(os.Stderr, "wrong number of integers on test %d\n", i+1)
			os.Exit(1)
		}
		got := make([]int, n)
		for j, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid integer in output on test %d\n", i+1)
				os.Exit(1)
			}
			got[j] = v
		}
		want := solve(n)
		for j := 0; j < n; j++ {
			if got[j] != want[j] {
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

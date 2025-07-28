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

func generateTests() [][2]int {
	r := rand.New(rand.NewSource(3))
	tests := make([][2]int, 100)
	for i := range tests {
		n := r.Intn(20) + 2
		x := r.Intn(n-1) + 2
		tests[i] = [2]int{n, x}
	}
	return tests
}

func buildFunnyPerm(n, x int) []int {
	if n%x != 0 {
		return nil
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	p[1] = x
	p[n] = 1
	cur := x
	for i := x + 1; i < n; i++ {
		if i%cur == 0 && n%i == 0 {
			p[cur], p[i] = p[i], p[cur]
			cur = i
		}
	}
	if cur != n {
		p[cur] = n
	}
	return p[1:]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc[0], tc[1])
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
	for i, tc := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test case %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "-1" {
			if buildFunnyPerm(tc[0], tc[1]) != nil {
				fmt.Fprintf(os.Stderr, "unexpected -1 on test %d\n", i+1)
				os.Exit(1)
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != tc[0] {
			fmt.Fprintf(os.Stderr, "wrong length on test %d\n", i+1)
			os.Exit(1)
		}
		want := buildFunnyPerm(tc[0], tc[1])
		if want == nil {
			fmt.Fprintf(os.Stderr, "expected -1 on test %d\n", i+1)
			os.Exit(1)
		}
		for j, p := range parts {
			v, err := strconv.Atoi(p)
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

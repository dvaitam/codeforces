package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n int
	s string
}

func expected(n int, s string) int {
	if n%2 == 1 {
		return -1
	}
	count := 0
	for _, ch := range s {
		if ch == '(' {
			count++
		} else {
			count--
		}
	}
	if count != 0 {
		return -1
	}
	balance := 0
	start := -1
	ans := 0
	for i, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 && start == -1 {
			start = i
		}
		if balance == 0 && start != -1 {
			ans += i - start + 1
			start = -1
		}
	}
	return ans
}

func run(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n := 0
		fmt.Sscan(parts[0], &n)
		s := parts[1]
		want := expected(n, s)
		got, err := run(bin, testCase{n: n, s: s})
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var g int
		fmt.Sscan(got, &g)
		if g != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, g)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

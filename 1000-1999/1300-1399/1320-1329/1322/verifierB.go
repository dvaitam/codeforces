package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	a []int
}

func expected(tc testCase) int {
	ans := 0
	for i := 0; i < tc.n; i++ {
		for j := i + 1; j < tc.n; j++ {
			ans ^= tc.a[i] + tc.a[j]
		}
	}
	return ans
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts)-1 != n {
			fmt.Fprintf(os.Stderr, "test %d expected %d numbers got %d\n", idx, n, len(parts)-1)
			os.Exit(1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], _ = strconv.Atoi(parts[1+i])
		}
		tc := testCase{n: n, a: a}
		want := expected(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(gotStr)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

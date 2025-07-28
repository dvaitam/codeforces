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

type testCaseD struct {
	n int
	a []int
	b []int
}

func parseTests(path string) ([]testCaseD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseD, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &a[j], &b[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseD{n: n, a: a, b: b}
	}
	return cases, nil
}

func solve(tc testCaseD) int {
	n := tc.n
	a := tc.a
	b := tc.b
	if n > 16 {
		return -1
	}
	m := 1 << n
	minFlip := -1
	for mask := 0; mask < m; mask++ {
		front := make([]int, n)
		back := make([]int, n)
		flips := 0
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				front[i] = b[i]
				back[i] = a[i]
				flips++
			} else {
				front[i] = a[i]
				back[i] = b[i]
			}
		}
		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if front[idx[i]] > front[idx[j]] {
					idx[i], idx[j] = idx[j], idx[i]
				}
			}
		}
		ok := true
		for i := 1; i < n; i++ {
			if back[idx[i-1]] <= back[idx[i]] {
				ok = false
				break
			}
		}
		if ok {
			if minFlip == -1 || flips < minFlip {
				minFlip = flips
			}
		}
	}
	return minFlip
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expect := solve(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
		}
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || val != expect {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

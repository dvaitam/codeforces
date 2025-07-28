package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n   int
	arr []int
}

func parseTests(path string) ([]testCaseB, error) {
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
	cases := make([]testCaseB, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		total := n * n
		arr := make([]int, total)
		for j := 0; j < total; j++ {
			if _, err := fmt.Fscan(in, &arr[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseB{n: n, arr: arr}
	}
	return cases, nil
}

type cell struct{ x, y int }

func solve(tc testCaseB) []string {
	n := tc.n
	even := make([]cell, 0)
	odd := make([]cell, 0)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if (i+j)%2 == 0 {
				even = append(even, cell{i, j})
			} else {
				odd = append(odd, cell{i, j})
			}
		}
	}
	idxE, idxO := 0, 0
	res := make([]string, 0, n*n)
	for _, a := range tc.arr {
		if idxE < len(even) && idxO < len(odd) {
			if a != 1 {
				c := even[idxE]
				idxE++
				res = append(res, fmt.Sprintf("1 %d %d", c.x, c.y))
			} else {
				c := odd[idxO]
				idxO++
				res = append(res, fmt.Sprintf("2 %d %d", c.x, c.y))
			}
		} else if idxE < len(even) {
			c := even[idxE]
			idxE++
			if a == 1 {
				res = append(res, fmt.Sprintf("3 %d %d", c.x, c.y))
			} else {
				res = append(res, fmt.Sprintf("1 %d %d", c.x, c.y))
			}
		} else {
			c := odd[idxO]
			idxO++
			if a == 2 {
				res = append(res, fmt.Sprintf("3 %d %d", c.x, c.y))
			} else {
				res = append(res, fmt.Sprintf("2 %d %d", c.x, c.y))
			}
		}
	}
	return res
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expect := solve(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != len(expect) {
			fmt.Printf("case %d failed: expected %d lines got %d\n", idx+1, len(expect), len(lines))
			os.Exit(1)
		}
		for i := range lines {
			if strings.TrimSpace(lines[i]) != expect[i] {
				fmt.Printf("case %d line %d mismatch\nexpected: %s\n got: %s\n", idx+1, i+1, expect[i], lines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

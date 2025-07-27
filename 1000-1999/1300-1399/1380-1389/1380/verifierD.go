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
	n, m    int
	x, k, y int64
	a, b    []int64
}

func parseTestcases(path string) ([]testCaseD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseD, T)
	for i := 0; i < T; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var x, k, y int64
		fmt.Fscan(in, &x, &k, &y)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
		}
		b := make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &b[j])
		}
		cases[i] = testCaseD{n: n, m: m, x: x, k: k, y: y, a: a, b: b}
	}
	return cases, nil
}

func getBound(a []int64, idx int) int64 {
	if idx >= 0 && idx < len(a) {
		return a[idx]
	}
	return 0
}

func processSegment(seg []int64, left, right, x, k, y int64) (int64, bool) {
	l := int64(len(seg))
	if l == 0 {
		return 0, true
	}
	maxV := left
	if right > maxV {
		maxV = right
	}
	maxInside := seg[0]
	for _, v := range seg {
		if v > maxInside {
			maxInside = v
		}
	}
	needFire := maxInside > maxV
	if needFire && l < k {
		return 0, false
	}
	var res int64
	if needFire {
		res += x
		l -= k
	}
	if k*y < x {
		res += l * y
	} else {
		res += (l/k)*x + (l%k)*y
	}
	if !needFire {
		berserkOnly := int64(len(seg)) * y
		if res > berserkOnly {
			res = berserkOnly
		}
	}
	return res, true
}

func solveCase(tc testCaseD) (int64, bool) {
	n, m := tc.n, tc.m
	a, b := tc.a, tc.b
	ai, bi := 0, 0
	last := -1
	var total int64
	for bi < m {
		start := ai
		for ai < n && a[ai] != b[bi] {
			ai++
		}
		if ai == n {
			return 0, false
		}
		cost, ok := processSegment(a[start:ai], getBound(a, last), a[ai], tc.x, tc.k, tc.y)
		if !ok {
			return 0, false
		}
		total += cost
		last = ai
		ai++
		bi++
	}
	cost, ok := processSegment(a[ai:], getBound(a, last), 0, tc.x, tc.k, tc.y)
	if !ok {
		return 0, false
	}
	total += cost
	return total, true
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.x, tc.k, tc.y))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected, ok := solveCase(tc)
		if !ok {
			if strings.TrimSpace(outStr) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d expected -1 got %s\n", idx+1, strings.TrimSpace(outStr))
				os.Exit(1)
			}
			continue
		}
		got, err := strconv.ParseInt(strings.TrimSpace(outStr), 10, 64)
		if err != nil || got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, expected, strings.TrimSpace(outStr))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

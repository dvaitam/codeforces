package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n   int
	arr []int
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
100
3
3 3 3 3 1 2 1 1 1
1
1
2
1 2 1 1
2
3 2 3 2
2
2 3 2 1
1
3
2
1 3 3 3
1
1
3
1 3 1 2 2 2 2 1 3
3
2 1 2 3 2 2 3 2 2
2
2 1 2 3
3
1 3 3 1 1 2 3 3 1
2
1 2 3 2
3
2 1 3 1 3 1 3 2 3
1
1
1
2
2
3 1 2 2
3
1 2 3 1 2 2 2 2 1
2
1 1 2 1
2
2 3 1 1
1
3
1
2
1
2
2
3 1 3 2
2
1 1 2 1
3
2 2 2 1 1 2 3 1 1
2
3 2 3 1
2
3 1 2 1
1
1
2
3 3 1 1
3
2 3 2 2 2 2 3 1 3
3
2 2 3 3 1 2 1 2 2
1
3
3
1 3 2 3 3 1 2 2 3
2
2 2 2 3
3
3 3 1 3 2 1 3 2 3
1
2
3
1 3 3 2 3 3 2 3 2
2
2 2 1 1
2
3 1 2 1
1
1
2
3 3 1 2
3
1 2 2 1 2 2 3 3 2
3
2 3 2 3 1 1 1 2 1
1
2
1
2
2
3 2 3 1
3
2 1 2 2 3 1 3 2 3
1
1
1
3
1
2
1
2
1
1
3
3 1 2 1 3 2 2 2 3
3
2 3 2 3 3 1 1 3 3
1
3
1
3
2
2 3 1 3
2
1 3 2 2
1
3
1
2
3
3 1 2 3 2 2 3 3 1
2
2 3 1 2
1
1
1
1
1
1
1
3
2
3 1 2 1
3
3 3 1 1 2 2 1 2 3
1
2
1
2
1
3
2
3 1 1 3
1
2
2
2 3 2 3
2
1 2 1 2
1
1
1
2
1
1
2
2 2 2 3
2
1 3 1 2
1
1
3
3 2 1 3 1 3 3 2 1
1
2
3
3 2 3 1 2 2 1 2 3
1
1
1
2
3
3 3 1 2 1 1 1 3 1
2
2 3 2 2
3
2 3 1 2 1 3 2 2 1
2
2 2 3 2
2
2 1 3 3
1
1
1
1
1
2
1
3
3
1 2 2 3 3 3 2 3 2
3
2 2 3 3 1 2 2 3 2
3
2 1 3 1 3 3 2 1 3
1
3
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		total := n * n
		arr := make([]int, total)
		for j := 0; j < total; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of test data")
	}
	return res, nil
}

type cell struct{ x, y int }

// solve mirrors 1503B.go for one test case.
func solve(tc testCase) []string {
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

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expect := solve(tc)
		out, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d: %v", idx+1, err)
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

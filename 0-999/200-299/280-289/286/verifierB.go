package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt.
const testcasesBData = `
4
41
44
18
32
6
7
45
50
10
58
11
4
55
7
59
46
55
36
45
27
55
47
35
19
35
53
17
56
15
59
45
39
54
28
39
19
30
33
44
43
46
60
52
24
7
22
41
9
33
39
42
23
56
14
17
3
48
19
9
47
16
25
52
12
23
29
54
5
8
52
11
56
46
16
4
54
38
42
60
36
40
45
6
3
9
42
14
40
55
38
9
27
7
25
55
9
4
40
3
`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// solve mirrors 286B.go, returning the final permutation.
func solve(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	for k := 2; k <= n; k++ {
		for s := 0; s < n; s += k {
			L := min(k, n-s)
			if L <= 1 {
				continue
			}
			tmp := p[s]
			copy(p[s:s+L-1], p[s+1:s+L])
			p[s+L-1] = tmp
		}
	}
	return p
}

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesBData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n})
	}
	return cases, nil
}

func runCandidate(bin string, n int) ([]int, error) {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("parse output %q: %v", f, err)
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.n)
		got, err := runCandidate(bin, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		for j := 0; j < tc.n; j++ {
			if got[j] != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d: expected %d got %d\n", i+1, j+1, expect[j], got[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

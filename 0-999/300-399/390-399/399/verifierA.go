package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt.
const testcasesAData = `
52 49 27
8 5 8
54 51 54
41 31 23
77 28 65
20 10 5
99 13 80
35 35 10
42 7 5
90 43 61
74 13 46
58 21 40
84 27 71
64 57 34
10 9 1
14 12 14
54 46 53
88 81 1
81 64 43
34 21 5
27 19 8
33 10 29
14 2 6
68 63 14
41 36 19
93 16 71
45 35 14
80 71 76
39 29 6
79 50 41
76 31 38
26 7 6
7 5 6
36 31 5
14 11 13
19 5 2
13 12 9
90 51 68
38 34 16
30 29 22
78 54 75
38 29 32
87 83 46
13 6 10
17 16 11
27 8 1
96 35 15
93 29 48
24 11 14
10 2 3
92 29 6
76 69 10
6 1 6
27 20 27
76 16 51
14 6 14
17 2 1
27 6 23
18 16 7
96 8 87
5 5 4
82 13 34
11 4 2
85 39 45
58 12 4
67 60 6
79 13 51
28 9 12
96 61 73
24 23 22
29 25 2
89 21 21
46 34 17
18 15 6
4 4 4
75 66 40
86 46 50
87 33 20
74 2 59
97 11 43
97 6 70
38 9 16
100 62 46
81 37 46
78 17 40
52 48 27
86 11 1
79 25 43
23 8 8
84 58 49
93 87 73
56 3 26
92 73 54
87 6 22
60 5 17
92 21 58
70 63 1
7 4 3
42 30 4
56 13 36
`

// solve mirrors 399A.go.
func solve(n, p, k int) string {
	start := p - k
	if start < 1 {
		start = 1
	}
	end := p + k
	if end > n {
		end = n
	}
	var tokens []string
	if start > 1 {
		tokens = append(tokens, "<<")
	}
	for i := start; i <= end; i++ {
		if i == p {
			tokens = append(tokens, fmt.Sprintf("(%d)", i))
		} else {
			tokens = append(tokens, strconv.Itoa(i))
		}
	}
	if end < n {
		tokens = append(tokens, ">>")
	}
	return strings.Join(tokens, " ")
}

type testCase struct {
	n, p, k int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesAData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 integers", idx+1)
		}
		n, _ := strconv.Atoi(fields[0])
		p, _ := strconv.Atoi(fields[1])
		k, _ := strconv.Atoi(fields[2])
		cases = append(cases, testCase{n: n, p: p, k: k})
	}
	return cases, nil
}

func runBinary(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("%d %d %d\n", tc.n, tc.p, tc.k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.n, tc.p, tc.k)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n   got: %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

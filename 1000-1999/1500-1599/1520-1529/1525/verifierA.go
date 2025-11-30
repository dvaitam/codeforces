package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// gcd mirrors 1525A.go.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	k int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
97
98
99
100
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		k, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("case %d bad integer: %v", i+1, err)
		}
		res = append(res, testCase{k: k})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1525A.go.
func solve(tc testCase) string {
	g := gcd(tc.k, 100)
	return fmt.Sprintf("%d", 100/g)
}

func runCandidate(bin string, tcs []testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tcs)))
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d\n", tc.k))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	got, err := runCandidate(bin, tests)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outputs := strings.Fields(got)
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d outputs got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc)
		if outputs[i] != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expect, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

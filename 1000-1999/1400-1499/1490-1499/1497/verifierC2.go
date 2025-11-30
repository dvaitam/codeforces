package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// referenceSolve mirrors 1497C2.go for a single test case.
func referenceSolve(n, k int64) string {
	d := k - 3
	n2 := n - d
	var parts []int64
	if n2%3 == 0 {
		a := n2 / 3
		parts = []int64{a, a, a}
	} else if n2%2 == 0 {
		half := n2 / 2
		if half%2 == 1 {
			x := (n2 - 2) / 2
			parts = []int64{x, x, 2}
		} else {
			parts = []int64{half, half / 2, half / 2}
		}
	} else {
		x := (n2 - 1) / 2
		parts = []int64{x, x, 1}
	}
	for i := int64(0); i < d; i++ {
		parts = append(parts, 1)
	}
	out := make([]string, len(parts))
	for i, v := range parts {
		out[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(out, " ")
}

type testCase struct {
	n int64
	k int64
}

// Embedded testcases from testcasesC2.txt.
const testcaseData = `
63 40
142 36
97 80
124 83
151 19
158 6
123 36
144 62
52 48
123 72
143 124
104 84
41 17
165 41
136 102
192 6
174 19
43 40
13 7
10 7
124 79
187 102
185 112
104 96
150 116
37 26
27 4
37 34
58 19
175 114
163 80
110 67
101 76
92 71
152 107
152 62
89 6
74 23
181 86
141 29
185 170
57 43
149 71
75 18
19 18
166 126
25 14
20 16
41 4
78 57
199 109
33 4
157 14
99 94
153 87
144 74
132 63
12 7
4 3
30 22
140 11
53 29
77 36
42 5
89 43
95 20
99 51
120 114
136 101
167 155
177 146
29 22
132 72
113 84
187 186
63 62
80 58
69 69
80 73
89 4
109 77
83 5
99 81
153 37
18 13
122 48
176 93
158 74
191 128
8 7
18 3
97 35
163 119
79 78
156 84
48 26
50 23
197 97
155 70
79 51
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields", i+1)
		}
		n, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		k, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %v", i+1, err)
		}
		res = append(res, testCase{n: n, k: k})
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	input := sb.String()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected := referenceSolve(tc.n, tc.k)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput:\n1\n%d %d\nexpected: %s\ngot: %s\n", idx+1, tc.n, tc.k, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

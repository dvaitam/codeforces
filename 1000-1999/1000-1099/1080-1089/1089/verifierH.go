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

type testcase struct {
	a      int64
	b      int64
	expect int64
}

const testcasesRaw = `98 18 2
67 101 1
56 120 8
53 146 1
10 130 10
15 48 3
39 156 39
124 172 4
127 5 1
91 80 1
153 196 1
103 80 1
139 13 1
2 176 2
94 123 1
173 71 1
158 131 1
174 195 3
2 159 1
25 55 5
127 77 1
144 185 1
157 89 1
86 87 1
9 55 1
186 166 2
39 24 3
23 159 1
95 97 1
160 144 16
73 36 1
24 166 2
119 185 1
91 164 1
167 146 1
20 195 5
133 46 1
141 69 3
160 137 1
21 102 3
124 164 4
123 198 3
175 24 1
78 38 2
170 125 5
40 60 20
2 39 1
117 82 1
82 61 1
102 192 6
152 31 1
71 31 1
49 38 1
165 101 1
189 167 1
142 162 2
199 19 1
145 167 1
107 106 1
80 79 1
11 148 1
85 24 1
112 138 2
73 36 1
193 194 1
67 151 1
14 28 14
25 111 1
46 37 1
135 18 9
24 5 1
96 185 1
135 103 1
107 54 1
93 200 1
106 52 2
29 71 1
70 152 2
76 171 19
132 165 33
3 46 1
97 186 1
76 76 76
141 200 1
71 152 1
66 119 1
187 63 1
175 154 7
198 101 1
26 141 1
6 54 6
72 69 3
157 61 1
167 191 1
68 170 34
100 158 2
195 87 3
137 111 1
140 67 1
99 175 1`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	readInt := func() (int64, bool) {
		if !scanner.Scan() {
			return 0, false
		}
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		return v, true
	}

	var res []testcase
	for {
		a, ok := readInt()
		if !ok {
			break
		}
		b, ok := readInt()
		if !ok {
			panic("dangling value: expected b")
		}
		expect, ok := readInt()
		if !ok {
			panic("dangling value: expected expected result")
		}
		res = append(res, testcase{a: a, b: b, expect: expect})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(res) == 0 {
		panic("no testcases parsed")
	}
	return res
}

// gcd computes the greatest common divisor using the Euclidean algorithm.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// solve mirrors the logic in 1089H.go: normalize to non-negative then output gcd.
func solve(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return gcd(a, b)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	v, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return v, nil
}

func checkCase(bin string, idx int, tc testcase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
	expected := solve(tc.a, tc.b)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `100
41 4
24 15
42 2
37 6
20 8
35 17
49 6
13 19
38 13
18 9
30 17
25 15
33 0
32 12
27 14
41 1
14 9
21 1
28 2
17 5
6 2
28 9
8 14
46 10
9 18
24 4
22 10
23 0
2 9
16 11
27 11
27 14
20 0
38 18
43 14
35 4
46 14
28 11
49 5
8 10
14 17
45 5
16 16
21 1
36 1
9 19
45 17
24 11
40 4
48 18
2 5
12 5
2 4
26 4
49 18
16 4
12 16
45 5
13 12
24 15
11 19
32 12
14 5
42 2
22 4
5 3
11 17
29 12
25 10
49 7
43 9
35 17
40 12
24 7
31 2
40 4
2 14
26 7
30 18
19 4
27 6
40 5
29 2
41 15
39 6
18 12
34 3
34 7
49 0
43 2
10 15
19 9
48 4
28 0
38 3
26 14
9 8
9 19
3 18
22 12`

type testCase struct {
	n int64
	m int64
}

func rangeOr(l, r int64) int64 {
	res := r
	for i := 0; i < 61; i++ {
		if l>>i < r>>i {
			res |= 1 << i
		}
	}
	return res
}

func solveCase(n, m int64) int64 {
	l := n - m
	if l < 0 {
		l = 0
	}
	r := n + m
	return rangeOr(l, r)
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesB)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].n, &tests[i].m); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.FormatInt(tc.n, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(tc.m, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	allInput := buildAllInput(tests)
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}

	outLines := strings.Fields(allOutput)
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.n, tc.m)
		got, err := strconv.ParseInt(outLines[i], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %q\n", i+1, outLines[i])
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d\nexpected: %d\ngot: %d\n", i+1, tc.n, tc.m, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

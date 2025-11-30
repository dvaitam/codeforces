package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `100
6 0
8 4
7 6
10 9
5 9
6 13
2 4
3 4
9 4
13 9
9 2
11 2
1 10
3 13
8 1
10 4
5 9
11 10
8 7
12 11
4 0
8 0
1 6
6 11
10 0
14 3
5 3
10 5
1 3
5 4
3 2
7 5
1 5
5 8
1 4
5 6
1 8
3 12
3 9
7 13
4 7
4 11
6 5
10 6
4 2
5 3
0 9
5 11
7 1
7 6
2 0
2 5
8 10
11 15
8 4
12 5
3 10
7 13
9 4
12 7
10 5
10 7
9 1
12 5
10 5
11 6
0 4
0 9
3 5
4 7
6 0
6 1
3 0
7 5
8 9
13 9
0 1
5 2
9 1
12 1
5 1
6 5
0 3
1 8
1 7
2 12
0 10
0 14
6 9
6 11
1 3
1 8
4 5
7 6
0 8
3 8
9 1
14 4
3 4
5 9
7 9
8 14
10 3
10 8
2 5
6 7
1 9
4 14
2 0
6 5
6 9
10 11
10 5
13 10
4 2
8 7
0 7
5 7
5 0
9 2
2 3
6 5
9 4
14 6
9 10
13 11
4 6
10 9
10 1
10 5
3 5
4 6
3 10
6 13
10 9
13 9
6 9
9 14
0 2
3 2
4 2
7 6
7 8
11 8
0 7
2 9
7 0
10 1
8 10
8 15
2 0
6 5
6 5
7 6
0 10
4 14
1 3
1 7
10 3
12 5
2 1
5 4
10 1
10 3
7 1
9 2
10 8
15 13
5 1
6 3
0 3
6 5
8 5
10 9
0 9
5 12
10 7
15 10
5 8
6 9
6 9
8 9
2 4
4 6
5 1
7 5
0 4
1 5
9 4
11 7
8 2
10 2
7 3
7 5
2 8
7 8
4 6
6 8
6 1
6 5`

type testCase struct {
	x1, y1, x2, y2 int64
}

func solveCase(x1, y1, x2, y2 int64) string {
	diff1 := x1 - y1
	diff2 := x2 - y2
	if diff1*diff2 > 0 {
		return "YES"
	}
	return "NO"
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesA)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no tests")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if len(fields) != 1+t*4 {
		return nil, fmt.Errorf("expected %d numbers, got %d", 1+t*4, len(fields))
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		base := 1 + 4*i
		x1, _ := strconv.ParseInt(fields[base], 10, 64)
		y1, _ := strconv.ParseInt(fields[base+1], 10, 64)
		x2, _ := strconv.ParseInt(fields[base+2], 10, 64)
		y2, _ := strconv.ParseInt(fields[base+3], 10, 64)
		tests[i] = testCase{x1: x1, y1: y1, x2: x2, y2: y2}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n%d %d\n", tc.x1, tc.y1, tc.x2, tc.y2)
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
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
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
	outFields := strings.Fields(allOutput)
	if len(outFields) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.x1, tc.y1, tc.x2, tc.y2)
		if outFields[i] != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d %d %d\nexpected: %s\ngot: %s\n", i+1, tc.x1, tc.y1, tc.x2, tc.y2, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

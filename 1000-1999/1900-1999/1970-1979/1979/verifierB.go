package main

import (
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `100
8 36
48 4
16 7
31 48
28 30
41 24
50 13
6 31
1 24
27 38
48 49
0 44
28 17
46 14
37 6
20 1
1 41
34 0
24 43
13 27
46 1
33 14
48 28
31 35
14 22
14 43
14 48
29 18
1 26
35 41
6 11
40 46
18 7
47 21
46 45
32 27
32 42
12 19
18 37
31 32
25 37
2 30
15 47
25 26
42 11
23 35
44 49
43 47
23 5
28 42
32 6
49 10
33 25
23 31
46 1
30 2
19 45
39 37
37 25
41 10
10 32
14 0
49 12
34 35
14 25
32 22
36 22
29 17
42 35
38 46
0 24
50 47
32 8
33 49
35 13
27 3
30 23
36 35
12 32
26 31
22 26
22 0
34 39
50 39
21 29
38 1
14 40
11 35
37 11
5 35
16 2
43 4
5 1
28 0
48 17
15 17
7 39
11 22
18 4
10 16`

type testCase struct {
	x, y int
}

func solveCase(x, y int) int {
	d := x ^ y
	return 1 << bits.TrailingZeros(uint(d))
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesB)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return nil, err
		}
		tests[i] = testCase{x: x, y: y}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.x))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.y))
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
	outFields := strings.Fields(allOutput)
	if len(outFields) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.x, tc.y)
		got, err := strconv.Atoi(outFields[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %q\n", i+1, outFields[i])
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d\nexpected: %d\ngot: %d\n", i+1, tc.x, tc.y, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

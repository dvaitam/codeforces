package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `-2 1
-90 -1
30 2
3 -1
22 0
49 -2
29 -3
-28 -3
93 -4
58 -1
36 4
-63 -1
-75 -4
75 0
20 3
-75 0
11 0
56 5
-48 3
22 2
33 -1
-85 3
-97 -4
84 1
81 5
60 -5
56 2
-15 -2
86 0
80 -4
-52 4
-44 -2
-64 3
14 -4
-80 0
30 2
-73 -1
41 -1
80 -4
40 0
38 -2
54 3
50 -1
13 -4
52 1
-19 4
-39 -1
-53 -2
-53 -5
56 5
-34 2
-83 -4
73 -3
-62 -5
-80 3
74 1
80 3
-30 3
-40 -2
73 4
7 4
-30 2
26 5
64 0
-79 0
56 -4
24 4
61 0
-52 -2
-96 -1
-71 -2
-5 -3
-15 1
-85 -4
100 -3
78 -2
-89 4
62 3
54 5
-82 -5
-69 5
-52 4
47 -4
0 -4
-6 -4
-91 4
-95 -2
-53 -4
22 -2
86 -5
73 -5
39 1
58 -4
-34 -4
-44 -4
65 -1
-11 1
-54 -5
28 2
-90 4`

func expected(x, y int) string {
	if y >= -1 {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	x int
	y int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesA)
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("testcase data malformed")
	}
	tests := make([]testCase, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		x, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{x: x, y: y})
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
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
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(lines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := expected(tc.x, tc.y)
		if lines[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, lines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

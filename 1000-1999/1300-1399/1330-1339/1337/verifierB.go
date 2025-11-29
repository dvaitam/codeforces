package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesB.txt.
const testcasesRaw = `17612 18 27
8272 8 3
64938 24 14
61899 20 12
27520 3 15
3716 28 26
51094 13 19
99914 24 0
91205 14 8
94574 25 7
77484 30 3
41607 0 0
3336 20 17
1207 30 28
49966 21 6
55328 23 0
69158 7 24
57395 30 15
72465 7 11
30261 21 7
99739 14 30
37983 29 0
54550 26 29
72936 29 20
13108 5 20
94849 27 9
15846 23 10
94567 22 16
55327 16 26
87859 6 9
37246 18 28
65453 27 30
66229 12 18
4526 15 7
97483 25 12
54305 21 5
48120 17 28
92149 24 21
96760 11 2
57536 21 16
14147 24 5
68281 26 12
48566 15 23
3877 15 1
40440 22 27
80585 18 18
51590 20 5
22098 16 7
1613 24 6
70729 29 27
71872 7 12
67342 11 30
75733 11 14
35295 21 17
79816 30 23
749 12 25
97060 16 25
16941 16 24
73579 6 13
7357 15 27
47807 18 17
26194 30 16
54186 15 26
46766 13 11
208 17 17
81723 25 19
43403 14 19
3667 25 7
83280 5 17
76607 5 27
12007 25 17
33462 1 26
88227 2 2
2188 14 0
98848 24 8
32711 8 3
81895 5 11
38049 2 5
20923 8 16
22040 21 8
84962 22 9
59599 22 10
65077 15 3
3098 9 12
45003 13 25
24647 8 3
33222 28 23
66862 6 30
79384 13 26
2729 7 0
52077 4 1
94220 30 5
58415 22 16
88890 13 17
28915 20 25
91102 16 14
29255 16 20
4024 12 21
75478 25 10
86485 20 13
7706 23 9
16474 30 6
6219 9 2
10020 9 29
39044 23 5
54549 18 8
17091 0 17
4970 18 26
28521 30 28
74748 14 5
92278 19 16
4906 12 6
45473 3 6
75155 21 28
56748 18 6
64534 3 30
87289 12 9
66075 15 0
42644 19 27
52734 28 9`

type testCase struct {
	x int
	n int
	m int
}

// referenceSolution mirrors 1337B.go logic.
func referenceSolution(x, n, m int) bool {
	for n > 0 && x > 20 {
		x = x/2 + 10
		n--
	}
	x -= 10 * m
	return x <= 0
}

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			panic("invalid testcase line")
		}
		x, _ := strconv.Atoi(fields[0])
		n, _ := strconv.Atoi(fields[1])
		m, _ := strconv.Atoi(fields[2])
		tests = append(tests, testCase{x: x, n: n, m: m})
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		expected := "NO"
		if referenceSolution(tc.x, tc.n, tc.m) {
			expected = "YES"
		}
		input := fmt.Sprintf("1\n%d %d %d\n", tc.x, tc.n, tc.m)
		outStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.ToUpper(strings.TrimSpace(outStr))
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %q\n", idx+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

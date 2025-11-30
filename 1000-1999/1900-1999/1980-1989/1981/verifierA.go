package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `100
42 85
49 110
36 91
27 68
12 42
20 40
21 46
39 95
21 50
21 42
31 67
47 97
32 82
8 35
17 34
11 27
30 61
21 49
13 44
15 35
25 58
37 85
1 12
33 67
9 24
10 39
38 93
50 100
22 54
33 76
35 81
46 107
2 23
31 81
10 38
3 19
49 115
37 90
9 20
7 32
35 79
19 46
43 100
16 37
7 32
20 57
13 32
45 100
49 100
34 73
29 77
30 79
20 58
25 64
2 13
46 98
7 27
5 11
13 43
49 107
22 58
17 40
33 67
28 61
33 76
27 73
13 43
48 99
15 40
27 55
3 19
49 102
24 64
7 29
16 42
31 79
26 54
5 21
27 61
42 102
32 65
13 32
24 51
7 24
3 8
13 26
1 2
11 34
24 48
6 22
37 89
20 52
4 27
49 98
3 13
37 90
16 32
22 49
25 55
8 26`

func countFactors(x int) int {
	cnt := 0
	d := 2
	for d*d <= x {
		for x%d == 0 {
			cnt++
			x /= d
		}
		d++
	}
	if x > 1 {
		cnt++
	}
	return cnt
}

func solveCase(l, r int) int {
	best := 0
	for x := l; x <= r; x++ {
		if c := countFactors(x); c > best {
			best = c
		}
	}
	return best
}

type testCase struct {
	l, r int
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesA)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].l, &tests[i].r); err != nil {
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
		sb.WriteString(strconv.Itoa(tc.l))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.r))
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
		want := solveCase(tc.l, tc.r)
		got, err := strconv.Atoi(outFields[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %q\n", i+1, outFields[i])
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d\nexpected: %d\ngot: %d\n", i+1, tc.l, tc.r, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

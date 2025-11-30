package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `100
31 9 11
22 3 16
8 9 14
3 6 11
48 4 6
41 8 3
5 3 7
48 5 2
25 2 4
26 10 17
19 9 16
38 5 14
6 7 8
17 4 14
13 7 4
5 2 17
29 5 4
32 8 9
14 2 7
40 4 4
13 9 13
24 10 5
7 9 5
37 8 14
34 9 11
32 9 7
35 5 1
22 7 11
3 10 5
17 4 13
38 6 16
5 3 17
3 3 8
9 2 10
1 9 11
11 4 15
24 10 13
34 10 2
37 3 17
49 3 14
49 5 10
35 8 16
25 5 1
43 2 6
20 10 19
17 7 3
32 6 10
50 8 13
25 2 6
42 4 8
19 7 2
3 9 14
10 9 20
46 3 5
23 8 2
40 9 13
30 2 4
31 4 1
3 4 11
7 10 12
13 8 16
8 2 20
45 9 20
41 7 4
44 6 5
25 6 4
34 5 2
26 9 12
49 5 15
23 3 2
3 9 9
2 10 19
37 5 8
6 10 17
27 10 10
8 4 14
37 8 3
7 8 3
7 8 5
47 2 15
28 8 1
32 7 9
6 7 3
8 7 1
23 7 6
1 5 12
5 4 7
1 5 4
48 2 10
24 2 20
15 4 6
30 3 16
23 6 5
2 5 12
22 9 10
19 10 11
12 3 4
35 6 6
25 4 5
15 7 17`

type testCase struct {
	x int64
	y int64
	k int64
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

func solveCase(x, y, k int64) int64 {
	for k > 0 && x >= y {
		rem := int64(y-1) - x%y
		if rem >= k {
			x += k
			k = 0
			break
		}
		x += rem
		k -= rem
		x++
		k--
		for x%y == 0 {
			x /= y
		}
	}
	if k > 0 {
		period := int64(y - 1)
		if period > 0 {
			k %= period
			x = ((x - 1 + k) % period) + 1
		}
	}
	return rangeOr(x, x) // x already adjusted
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesB)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].x, &tests[i].y, &tests[i].k); err != nil {
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
		sb.WriteString(strconv.FormatInt(tc.x, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(tc.y, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(tc.k, 10))
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
		want := solveCase(tc.x, tc.y, tc.k)
		got, err := strconv.ParseInt(outFields[i], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %q\n", i+1, outFields[i])
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %d %d %d\nexpected: %d\ngot: %d\n", i+1, tc.x, tc.y, tc.k, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

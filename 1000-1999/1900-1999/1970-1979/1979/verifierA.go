package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxInt = int(^uint(0) >> 1)

func solveCase(a []int) int {
	best := maxInt
	for i := 0; i < len(a)-1; i++ {
		mx := a[i]
		if a[i+1] > mx {
			mx = a[i+1]
		}
		if mx < best {
			best = mx
		}
	}
	return best - 1
}

type testCase struct {
	arr []int
}

const testcasesA = `100
8
14 2 9 17 16 13 10 16
7
19 7 17 5 10 5 4
6
18 20 5 10 4 3
7
16 18 4 12 14 11 20
5
18 16 15 17 9
2
18 1
3
13 1 20
9
11 8 11 3 7 19 8 8 5
10
15 3 3 11 17 16 4 10 18 10
3
18 11 18
5
20 18 19 10 15
3
20 13 11
5
10 6 7 6 2
6
16 3 3 5 5 2
3
18 13 17
6
17 8 7 19 14 19
6
15 16 12 3 11 20
3
16 19 11
5
8 1 9 4 8
7
6 11 14 2 4 5 8
2
19 18
3
1 4 7
3
13 3 12
3
2 20 1
5
6 4 16 7 2
2
18 14
3
9 3 8
3
10 12 14
4
2 17 15 2
3
13 7 9
7
16 19 6 7 2 6 6
7
17 9 4 20 15 6 1
9
14 19 17 10 12 13 9 5 18
2
15 3
7
2 18 9 5 8 16 12
6
12 19 20 5 10 13
8
3 1 20 7 11 6 8 8
9
13 19 14 2 13 19 14 2 6
9
3 9 6 15 17 16 18 20 1
2
16 11
6
15 2 14 7 18 3
4
1 13 14 11
2
7 1
2
17 20
3
7 4 20
5
10 9 6 4 16
8
3 1 9 15 4 9 5 17
7
4 5 9 1 2 2 7
6
18 11 12 19 2 20
9
15 14 12 18 6 7 13 19 10
2
5 5
6
11 11 12 3 11 20
2
2 9
4
5 19 10 12
8
18 5 10 4 16 8 2 10
4
17 3 10 13
7
10 14 4 4 18 16 16
7
11 4 16 4 16 14 2
6
11 5 6 19 13 3
3
3 7 8
2
13 1
3
13 18 17
6
15 16 19 7 14 3
7
8 9 19 6 14 7 12
3
3 1 17
9
7 4 16 13 9 7 2 7 20
4
4 7 15 13
7
18 5 4 20 16 5 19
8
14 17 16 11 16 16 7 18
5
1 11 11 11 2
10
5 9 20 5 13 19 10 16 3 3
10
2 3 8 5 2 10 1 15 11 6
4
15 12 17 13
10
17 2 19 3 17 20 3 14 7 10
10
20 14 16 13 20 19 8 1 1 6
6
17 19 9 11 3 16
6
10 14 13 13 2 6
4
8 10 11 2
2
16 14
4
16 20 3 5
7
14 2 20 15 13 15 2
3
16 5 1
2
20 20
4
11 4 18 12
5
13 16 4 2 20
9
20 11 4 20 10 5 13 10 4
10
7 2 13 15 12 7 15 12 3 2
2
16 9
2
17 19
5
8 3 17 17 14
10
10 4 5 14 19 14 3 4 14 3
3
14 5 1
9
14 14 1 16 11 9 3 12 3
3
12 1 12
7
6 1 8 12 3 20 5
5
1 7 4 1 10
7
1 20 8 5 6 15 4
9
12 9 5 1 7 12 11 16 10
6
18 11 6 19 3 4
10
19 10 6 13 5 5 8 11 17 8`

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesA)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, err
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				return nil, err
			}
		}
		tests[i] = testCase{arr: arr}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
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
	outLines := strings.Fields(allOutput)
	if len(outLines) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outLines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.arr)
		got, _ := strconv.Atoi(outLines[i])
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %v\nexpected: %d\ngot: %d\n", i+1, tc.arr, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesF = `100
4
15 10 20 5
1
4 23
3
20 16 7
2
3 7
1 13
3
9 0 17
2
3 15
1 1
3
14 9 16
1
3 3
5
10 16 9 18 13
1
3 5
4
14 20 6 2
2
1 1
2 15
3
3 1 14
1
2 15
4
1 7 13 10
1
1 6
4
7 2 4 15
1
2 0
3
18 8 12
1
3 20
5
7 19 3 4 9
2
1 6
1 16
4
3 13 18 12
2
2 8
3 10
3
8 9 8
1
2 3
4
5 18 3 2
2
3 15
2 3
4
20 2 18 4
1
4 5
4
3 9 13 0
2
4 19
2 11
4
7 13 9 4
2
1 1
3 6
5
2 12 5 20 17
1
1 8
5
15 11 7 3 14
2
1 2
1 4
5
10 17 3 14 10
2
1 11
5 1
5
1 18 18 8 19
2
2 20
1 2
3
3 4 15
2
2 16
1 9
4
19 11 14 14
2
2 18
1 6
5
15 11 6 6 17
1
5 17
5
14 13 4 5 3
2
4 7
3 3
4
17 4 1 20
1
3 3
4
9 6 18 12
2
1 2
4 16
4
8 10 15 10
1
3 12
3
18 15 3
2
2 1
1 7
3
8 17 15
2
1 1
3 10
4
10 1 2 1
2
2 11
3 17
4
2 8 2 9
2
4 4
2 4
5
8 14 7 13 7
1
4 16
5
12 5 12 8 18
2
4 17
4 5
3
16 1 12
1
3 3
4
7 17 15 10
1
4 13
4
11 20 14 18
2
1 13
4 13
3
7 16 5
1
1 11
5
2 13 17 16 15
1
2 19
3
18 11 14
1
1 12
3
7 19 4
2
3 20
2 11
4
13 2 19 10
2
3 2
4 15
3
11 13 3
2
2 8
2 20
4
8 17 13 10
1
1 14
5
8 17 6 17 9
2
2 20
5 8
3
14 6 2
1
1 18
4
18 16 17 19
2
2 17
3 19
5
1 20 16 17 0
2
2 3
2 18
4
18 3 5 2
1
1 9
4
14 15 5 2
2
2 14
1 1
5
2 17 13 19 19
1
4 20
4
10 14 20 11
1
1 1
3
3 1 11
1
1 20
5
16 18 9 2 7
1
5 0
5
5 12 12 17 14
2
1 0
1 1
4
4 1 12 14
1
2 7
3
6 14 5
2
3 6
2 0
5
16 0 16 7 3
2
1 3
1 1
4
14 19 7 11
2
4 0
2 7
3
15 8 3
1
3 6
3
14 4 7
2
2 0
2 7
4
14 1 17 8
1
1 18
4
14 17 6 1
1
2 3
5
17 9 16 1 2
2
3 10
3 18
3
0 1 7
2
1 19
1 11
3
13 6 12
1
1 4
5
4 8 15 17 12
2
4 8
3 14
4
15 17 2 14
2
1 12
2 8
4
17 5 18 16
2
1 13
3 8
5
13 16 9 11 2
1
5 5
3
15 16 8
1
3 7
4
7 3 19 3
1
1 0
5
7 15 10 8 3
2
5 17
4 4
5
17 15 0 9 11
2
3 15
2 0
4
12 0 8 12
2
1 12
1 0
4
7 16 4 8
1
4 19
3
1 16 10
1
1 14
4
18 0 16 8
1
2 11
5
11 3 1 15 2
2
4 16
2 11
5
11 6 12 12 2
2
3 11
2 9
5
11 17 12 0 2
2
2 17
2 17
4
9 15 5 19
1
1 9
5
1 20 10 16 17
1
1 13
3
16 12 4
1
1 10
5
13 17 17 16 1
2
3 5
1 19
3
9 18 4
2
1 15
3 5
4
7 4 18 12
1
1 14
4
10 10 19 17
1
2 8
3
3 5 2
2
3 4
3 12
4
11 14 10 1
2
1 3
3 9
3
12 15 14
2
2 16
3 1
4
1 18 6 9
1
3 6
3
20 19 13
1
1 12
5
3 2 0 12 12
2
1 20
2 10
5
10 8 6 17 20
2
4 6
1 18
3
19 10 15
1
1 1
4
6 1 12 18
2
2 16
4 12
3
14 3 14
2
1 15
3 20
5
14 17 16 20 5
1
3 20
4
15 13 19 18
2
2 1
3 20
4
16 18 9 6
1
4 9
3
7 15 17
2
3 1
1 10`

func findSegment(a []int) (int, int) {
	n := len(a)
	l := 0
	for l+1 < n && a[l] <= a[l+1] {
		l++
	}
	if l == n-1 {
		return -1, -1
	}
	r := n - 1
	for r > 0 && a[r-1] <= a[r] {
		r--
	}
	minV, maxV := a[l], a[l]
	for i := l; i <= r; i++ {
		if a[i] < minV {
			minV = a[i]
		}
		if a[i] > maxV {
			maxV = a[i]
		}
	}
	for l > 0 && a[l-1] > minV {
		l--
	}
	for r < n-1 && a[r+1] < maxV {
		r++
	}
	return l + 1, r + 1
}

func solveCase(n int, arr []int, queries [][2]int) []string {
	res := make([]string, len(queries)+1)
	l, r := findSegment(arr)
	res[0] = fmt.Sprintf("%d %d", l, r)
	for i, q := range queries {
		arr[q[0]-1] = q[1]
		l, r = findSegment(arr)
		res[i+1] = fmt.Sprintf("%d %d", l, r)
	}
	return res
}

type testCase struct {
	n       int
	arr     []int
	queries [][2]int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesF)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j], err = nextInt()
			if err != nil {
				return nil, err
			}
		}
		q, err := nextInt()
		if err != nil {
			return nil, err
		}
		queries := make([][2]int, q)
		for j := 0; j < q; j++ {
			p, err := nextInt()
			if err != nil {
				return nil, err
			}
			v, err := nextInt()
			if err != nil {
				return nil, err
			}
			queries[j] = [2]int{p, v}
		}
		tests[i] = testCase{n: n, arr: arr, queries: queries}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(len(tc.queries)))
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(strconv.Itoa(q[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(q[1]))
			sb.WriteByte('\n')
		}
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
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	idx := 0
	for i, tc := range tests {
		expect := solveCase(tc.n, append([]int(nil), tc.arr...), tc.queries)
		for step, exp := range expect {
			if idx+1 >= len(outFields) {
				fmt.Printf("missing output for test %d step %d\n", i+1, step)
				os.Exit(1)
			}
			got := outFields[idx] + " " + outFields[idx+1]
			if got != exp {
				fmt.Printf("test %d step %d failed\nexpected: %s\ngot: %s\n", i+1, step, exp, got)
				os.Exit(1)
			}
			idx += 2
		}
	}
	if idx != len(outFields) {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

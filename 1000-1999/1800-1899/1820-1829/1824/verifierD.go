package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
1 4 1 4
2 3 4 4
2 2 3 3
3 4 4 4
4 4 3 3
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 2
2 1
1 2 2 2
1 2 2 2
4 2
4 3 4 2
4 4 1 1
2 2 4 4
3 5
3 3 2
1 2 1 2
2 3 2 3
3 3 1 1
3 3 3 3
2 2 3 3
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 3
3 3 3
3 3 2 3
3 3 3 3
3 3 2 3
5 4
2 5 1 2 4
2 3 1 4
1 3 5 5
2 4 1 5
1 2 5 5
3 3
2 2 3
2 3 1 1
2 3 2 3
2 3 1 1
3 1
1 3 3
2 2 1 2
1 2
1
1 1 1 1
1 1 1 1
2 3
1 1
1 2 2 2
1 2 2 2
2 2 1 2
3 1
1 2 1
1 1 3 3
2 1
2 2
2 2 2 2
1 2
1
1 1 1 1
1 1 1 1
4 5
4 3 4 3
3 3 4 4
1 1 2 3
1 2 3 4
1 1 3 4
2 4 1 3
3 3
3 1 3
3 3 3 3
2 2 3 3
3 3 3 3
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
5 5
2 2 5 1 4
4 5 4 5
1 3 5 5
1 5 4 5
5 5 1 5
4 5 1 5
2 2
2 2
2 2 1 2
1 1 1 2
2 4
1 2
1 2 1 2
2 2 1 2
2 2 2 2
1 1 2 2
2 2
2 2
1 2 1 2
2 2 2 2
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 2
2 1 4 1
4 4 3 4
1 3 3 4
5 3
4 1 2 4 4
2 3 4 4
4 5 3 5
1 1 5 5
2 2
2 1
1 1 1 1
1 2 1 2
5 4
5 2 1 5 3
2 5 2 2
1 4 1 2
1 1 1 5
4 5 2 4
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 5
3 2 3 4
2 3 3 4
4 4 2 3
3 3 1 4
2 2 1 4
3 3 1 3
4 3
2 1 2 1
2 2 4 4
3 3 2 2
3 3 4 4
3 1
2 3 3
2 3 3 3
1 2
1
1 1 1 1
1 1 1 1
4 2
3 4 2 4
3 4 2 4
2 3 4 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 2
3 1 2
1 2 3 3
2 2 1 1
5 2
4 3 1 3 5
3 3 4 4
4 4 2 4
1 1
1
1 1 1 1
2 4
1 1
2 2 2 2
2 2 1 1
2 2 1 2
2 2 1 2
5 2
3 2 5 1 1
4 4 2 3
4 5 3 4
3 4
1 3 1
1 1 1 3
3 3 2 2
2 2 1 2
3 3 2 2
1 2
1
1 1 1 1
1 1 1 1
3 5
2 1 3
3 3 3 3
3 3 3 3
3 3 2 3
2 3 1 2
2 2 1 2
4 1
4 4 4 4
4 4 4 4
5 4
3 5 3 3 4
1 1 5 5
4 5 2 5
2 2 1 4
2 3 2 5
2 4
1 2
1 2 1 1
2 2 1 2
2 2 2 2
1 1 1 1
5 5
2 3 2 1 3
2 2 2 2
3 4 4 5
2 3 5 5
3 3 5 5
3 4 5 5
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
4 3 3 2
4 4 1 1
4 4 3 3
2 2 1 2
1 2 2 3
1 2
1
1 1 1 1
1 1 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
4 3 2 2
3 4 2 4
3 3 1 2
2 3 1 4
4 4 1 2
4 3
1 1 3 2
1 3 2 3
3 3 3 3
3 4 4 4
4 4
2 3 4 4
4 4 1 4
1 3 3 4
3 4 2 2
4 4 3 3
5 2
2 2 2 5 4
1 1 4 4
3 4 3 4
1 2
1
1 1 1 1
1 1 1 1
5 3
2 5 2 3 3
5 5 2 3
4 4 3 5
1 2 5 5
2 2
2 2
2 2 1 1
2 2 2 2
1 2
1
1 1 1 1
1 1 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 3
1 1 3
2 2 2 2
2 3 3 3
3 3 3 3
4 3
4 1 1 3
3 3 2 3
3 4 2 4
4 4 2 4
2 5
1 1
2 2 2 2
2 2 1 2
1 1 1 1
2 2 1 2
1 2 1 1
5 5
4 5 3 2 2
3 4 2 4
3 4 2 3
2 3 5 5
4 5 1 2
5 5 1 2
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 5
2 1
1 2 2 2
2 2 2 2
1 2 2 2
1 1 2 2
2 2 1 1
5 1
1 1 1 2 5
5 5 4 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 2
3 1 1
2 2 3 3
1 2 1 3
5 1
3 3 3 4 2
2 4 3 5
5 5
3 1 4 1 1
2 3 4 4
4 5 5 5
2 4 4 4
5 5 4 5
1 2 4 5
4 1
1 2 3 3
2 3 1 1
2 2
2 2
2 2 1 2
1 1 2 2
4 3
3 1 2 4
2 2 2 3
2 3 4 4
3 4 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 3
1 1
1 1 1 2
1 2 1 2
2 2 1 1
2 4
2 2
1 1 2 2
2 2 1 1
2 2 1 1
1 2 1 2
3 4
1 3 3
2 3 1 2
2 3 2 3
3 3 2 3
3 3 3 3
4 4
4 2 2 2
2 2 2 3
1 2 1 2
1 1 3 4
1 3 4 4
5 5
2 1 3 5 4
5 5 1 3
5 5 2 3
5 5 3 5
1 1 2 2
3 4 3 4
2 1
1 1
2 2 2 2
2 4
2 1
2 2 2 2
1 1 2 2
2 2 2 2
2 2 1 2
4 5
2 4 3 2
2 3 1 2
3 3 2 2
4 4 3 3
1 3 4 4
3 3 2 4
2 4
1 1
1 2 2 2
1 2 2 2
1 2 2 2
1 1 2 2
5 2
1 1 5 5 3
3 3 2 2
2 5 1 5
1 1
1
1 1 1 1
4 2
2 2 4 2
3 3 4 4
3 4 1 2
4 3
4 3 4 2
1 2 3 4
1 1 2 2
4 4 2 2
4 5
3 2 2 4
4 4 2 3
1 1 2 2
1 4 1 4
2 4 1 1
2 2 3 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 1
1 3 2
2 2 1 3
3 3
3 2 3
3 3 2 3
2 3 2 2
3 3 2 2
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
2 2
1 1
2 2 2 2
1 2 1 2
3 1
1 2 3
3 3 2 2
5 1
5 3 4 1 4
1 2 4 4
2 5
1 2
1 2 2 2
2 2 1 2
2 2 1 2
2 2 1 2
1 2 1 2
2 3
1 1
2 2 1 2
2 2 1 2
2 2 2 2
4 3
1 2 3 2
1 1 1 1
1 3 3 4
2 3 3 3
`

type query struct {
	l int
	r int
	x int
	y int
}

type testCase struct {
	input    string
	expected string
}

func g(a []int, i, j int) int {
	if i > j {
		return 0
	}
	required := make(map[int]struct{})
	for p := i; p <= j; p++ {
		required[a[p-1]] = struct{}{}
	}
	for x := j; x >= 1; x-- {
		if _, ok := required[a[x-1]]; ok {
			delete(required, a[x-1])
			if len(required) == 0 {
				return x
			}
		}
	}
	return 0
}

func solveQueries(n int, arr []int, qs []query) []int {
	results := make([]int, len(qs))
	for idx, q := range qs {
		ans := 0
		for i := q.l; i <= q.r; i++ {
			for j := q.x; j <= q.y; j++ {
				if i <= j {
					ans += g(arr, i, j)
				}
			}
		}
		results[idx] = ans
	}
	return results
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	idx := 0
	nextInt := func() (int, error) {
		if idx >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[idx])
		idx++
		return v, err
	}

	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}

	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		qCount, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: bad q: %w", caseNum+1, err)
		}

		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad array value: %w", caseNum+1, err)
			}
			arr[i] = v
		}

		qs := make([]query, qCount)
		for i := 0; i < qCount; i++ {
			l, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query l: %w", caseNum+1, err)
			}
			r, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query r: %w", caseNum+1, err)
			}
			x, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query x: %w", caseNum+1, err)
			}
			y, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query y: %w", caseNum+1, err)
			}
			qs[i] = query{l: l, r: r, x: x, y: y}
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, qCount)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range qs {
			fmt.Fprintf(&sb, "%d %d %d %d\n", qu.l, qu.r, qu.x, qu.y)
		}

		results := solveQueries(n, arr, qs)
		var expSb strings.Builder
		for i, v := range results {
			if i > 0 {
				expSb.WriteByte('\n')
			}
			expSb.WriteString(strconv.Itoa(v))
		}

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: expSb.String(),
		})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runCandidate(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

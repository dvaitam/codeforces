package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n, m int
	a    [][]int
	b    [][]int
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
100
2 1
4
4
4
4
1 3
4 4 5
4 5 5
3 2
5 1
5 5
4 2
4 2
5 1
5 5
2 3
1 5 1
2 4 1
1 5 1
2 4 1
2 2
1 4
5 4
5 4
1 5
1 3
3 5 1
3 5 1
2 1
3
4
4
2
3 3
3 3 4
1 3 5
4 1 4
3 3 4
4 1 4
1 3 5
2 1
2
1
2
0
2 1
1
4
5
1
1 1
5
4
3 3
3 3 1
4 3 2
3 5 3
4 3 2
3 3 1
3 5 3
3 3
5 3 2
4 3 4
4 5 3
5 3 2
4 5 3
4 3 4
3 1
2
1
5
1
2
5
3 3
2 5 5
5 1 3
5 2 4
5 1 3
5 2 5
2 5 5
2 2
1 4
3 5
1 4
3 5
1 1
5
5
1 3
2 1 2
2 1 2
3 2
5 2
4 4
3 1
4 4
5 2
3 1
1 2
5 1
5 1
1 3
5 4 3
4 4 3
2 3
1 3 5
4 5 4
1 3 4
4 5 4
2 1
3
5
5
3
1 2
4 5
4 5
3 1
3
4
2
3
4
3
3 1
2
5
5
5
4
2
3 3
1 4 5
3 3 1
1 2 2
1 5 5
1 2 2
3 3 1
2 3
3 1 1
3 5 4
3 5 4
3 1 1
2 1
2
3
3
2
2 3
3 3 4
5 2 1
3 3 4
5 2 1
2 1
3
3
3
2
1 2
3 1
3 1
2 1
2
3
2
2
2 3
4 5 5
2 4 5
4 5 5
2 4 5
2 3
1 1 4
3 4 2
1 1 4
3 4 2
3 1
4
1
3
2
1
4
1 1
1
0
1 3
5 2 1
5 2 1
1 1
5
5
2 1
5
2
2
4
2 1
4
2
4
2
1 3
3 1 5
3 1 5
2 2
4 2
5 5
4 5
4 2
2 1
2
5
5
2
3 3
3 2 4
5 5 1
5 1 4
5 1 4
5 5 1
3 2 4
2 3
3 2 4
3 4 1
3 2 4
3 4 1
1 2
2 3
2 3
2 1
4
1
4
1
2 1
4
2
4
2
1 1
4
5
3 3
1 2 3
5 1 1
2 3 1
2 3 1
1 2 3
5 1 1
2 2
2 4
4 2
4 2
2 4
2 3
5 2 3
2 5 5
5 2 3
2 5 4
3 1
5
4
1
1
4
5
2 3
2 1 1
1 2 2
2 1 1
1 2 2
3 1
1
5
4
1
5
4
2 2
2 2
3 5
2 2
3 5
2 2
3 4
1 2
1 2
3 4
1 2
4 5
5 5
1 2
2 1
2 1
3 3
4 5 2
1 2 5
4 3 4
1 2 5
4 5 2
4 2 4
1 1
4
4
2 2
1 1
2 5
2 5
1 1
3 2
2 4
3 1
5 4
5 4
2 4
3 1
3 3
3 2 5
3 1 1
2 3 5
3 2 5
3 1 1
2 2 5
2 2
4 1
1 4
1 4
4 1
1 1
4
5
2 1
3
3
3
3
3 1
5
5
1
1
4
5
3 2
3 2
3 1
5 5
4 5
3 2
3 1
3 1
4
1
5
5
1
5
1 3
2 4 2
2 4 2
3 2
3 2
3 3
5 2
5 2
3 3
3 2
1 2
5 1
5 0
3 2
1 1
5 1
3 1
3 1
5 1
1 1
3 3
3 4 1
4 5 5
4 1 1
3 4 1
4 5 5
4 1 1
1 2
2 4
2 4
2 1
1
1
1
0
3 3
3 1 2
5 2 3
5 1 4
3 1 2
5 1 4
5 2 3
2 2
5 4
4 1
5 4
4 1
3 1
2
4
4
5
4
2
1 2
1 3
1 3
1 3
2 4 3
2 5 3
3 1
2
5
3
5
2
2
1 3
3 5 1
3 5 0
1 1
5
4
3 2
1 4
2 1
4 2
1 4
4 2
2 0
1 1
1
1
1 1
2
2
2 3
3 4 3
4 3 1
4 3 1
3 4 3
3 1
1
1
5
5
1
0
3 3
5 3 1
5 2 3
3 3 3
5 3 1
5 2 2
3 3 3
1 3
4 3 4
4 3 4
1 2
3 3
3 3
3 3
4 4 5
4 5 3
5 2 2
5 2 3
4 4 5
4 5 3
1 3
5 5 5
5 5 5
3 2
3 3
4 3
4 5
4 3
4 5
3 3
2 3
1 1 2
5 1 3
5 1 3
1 1 2
1 1
5
4
1 3
1 1 5
1 1 5
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad m: %v", i+1, err)
		}
		a := make([][]int, n)
		b := make([][]int, n)
		for r := 0; r < n; r++ {
			a[r] = make([]int, m)
			for c := 0; c < m; c++ {
				val, err := nextInt()
				if err != nil {
					return nil, fmt.Errorf("case %d A(%d,%d): %v", i+1, r, c, err)
				}
				a[r][c] = val
			}
		}
		for r := 0; r < n; r++ {
			b[r] = make([]int, m)
			for c := 0; c < m; c++ {
				val, err := nextInt()
				if err != nil {
					return nil, fmt.Errorf("case %d B(%d,%d): %v", i+1, r, c, err)
				}
				b[r][c] = val
			}
		}
		res = append(res, testCase{n: n, m: m, a: a, b: b})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens in testcase data")
	}
	return res, nil
}

// solveColumns mirrors 1500C logic: returns columns (0-indexed) that sort A to B.
func solveColumns(a, b [][]int) ([]int, bool) {
	n := len(a)
	m := len(a[0])
	inv := make([]int, m)
	for i := 1; i < n; i++ {
		for j := 0; j < m; j++ {
			if b[i][j] < b[i-1][j] {
				inv[j]++
			}
		}
	}
	cut := make([]bool, n-1)
	var qu []int
	for j := 0; j < m; j++ {
		if inv[j] == 0 {
			qu = append(qu, j)
		}
	}
	var ord []int
	for len(qu) > 0 {
		v := qu[len(qu)-1]
		qu = qu[:len(qu)-1]
		ord = append(ord, v)
		for i := 1; i < n; i++ {
			if !cut[i-1] && b[i-1][v] < b[i][v] {
				cut[i-1] = true
				for j := 0; j < m; j++ {
					if b[i-1][j] > b[i][j] {
						inv[j]--
						if inv[j] == 0 {
							qu = append(qu, j)
						}
					}
				}
			}
		}
	}
	aa := make([][]int, n)
	for i := 0; i < n; i++ {
		aa[i] = append([]int(nil), a[i]...)
	}
	for i := len(ord) - 1; i >= 0; i-- {
		col := ord[i]
		sort.SliceStable(aa, func(x, y int) bool { return aa[x][col] < aa[y][col] })
	}
	if !matEqual(aa, b) {
		return nil, false
	}
	return ord, true
}

func matEqual(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", tc.a[i][j])
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", tc.b[i][j])
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validateOutput(tc testCase, output string) bool {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return false
	}
	if fields[0] == "-1" {
		return len(fields) == 1
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil || len(fields) != k+1 {
		return false
	}
	cols := make([]int, k)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil || v < 1 || v > tc.m {
			return false
		}
		cols[i] = v - 1
	}
	aa := make([][]int, tc.n)
	for i := 0; i < tc.n; i++ {
		aa[i] = append([]int(nil), tc.a[i]...)
	}
	for i := k - 1; i >= 0; i-- {
		col := cols[i]
		sort.SliceStable(aa, func(x, y int) bool { return aa[x][col] < aa[y][col] })
	}
	return matEqual(aa, tc.b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		_, possible := solveColumns(tc.a, tc.b)
		out, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !validateOutput(tc, out) {
			fmt.Printf("case %d failed: invalid transformation\n", idx+1)
			os.Exit(1)
		}
		if possible && strings.Fields(out)[0] == "-1" {
			fmt.Printf("case %d failed: solution exists but candidate returned -1\n", idx+1)
			os.Exit(1)
		}
		if !possible && strings.Fields(out)[0] != "-1" {
			fmt.Printf("case %d failed: expected -1 but got %s\n", idx+1, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt to avoid filesystem dependency.
const testcasesDData = `
8 5 3 8 10 1 9 7 1 1 2
5 2 7 4 5 7 1 4
6 3 7 1 10 7 10 1 3
6 9 6 8 8 5 6 2 2 2
4 8 6 1 6 2 3 2
1 2 2 1 1
7 10 6 9 3 1 8 7 2 7 4
6 1 4 10 4 2 8 1 6
2 9 9 2 2 1
2 7 1 3 2 2 1
4 10 5 7 8 3 3 1 1
2 6 7 3 2 2 1
5 8 7 9 1 7 3 3 4 2
7 7 7 10 10 5 8 9 2 3 4
3 6 1 7 2 2 2
2 4 6 1 2
1 10 2 1 1
8 2 2 3 10 1 4 3 1 3 7 7 3
5 7 8 5 9 1 1 4
3 2 2 8 3 2 3 2
5 7 1 10 5 6 2 1 1
5 4 8 5 9 4 2 1 4
1 1 1 1
1 5 2 1 1
1 10 1 1
3 8 2 10 1 1
3 2 2 2 3 3 1 1
2 8 7 1 1
8 8 9 3 9 2 3 6 7 3 3 2 1
4 2 7 6 9 2 4 3
4 10 2 3 10 2 3 3
6 9 6 1 1 1 3 2 6 6
4 3 3 2 8 3 4 2 3
5 3 1 3 6 6 1 1
8 3 5 6 4 2 10 4 6 2 2 5
7 9 6 1 6 4 7 1 2 5 3
3 5 1 5 3 1 2 3
7 1 5 4 3 5 6 9 3 1 2 4
4 10 10 10 2 2 3 1
5 1 1 10 3 6 1 4
4 10 9 4 1 3 2 4 4
6 6 6 2 1 7 6 1 4
4 7 10 1 4 3 4 4 4
3 3 7 7 2 3 3
4 5 2 3 10 1 3
3 2 3 5 1 2
5 9 7 6 9 7 2 4 1
5 10 6 5 5 5 2 5 5
1 8 3 1 1 1
7 1 1 7 9 2 5 1 1 4
7 9 6 9 8 9 8 9 3 6 6 4
3 10 1 10 2 1 2
7 2 4 2 3 10 8 3 2 3 7
8 3 8 5 8 6 9 4 1 2 5 7
2 5 3 3 1 1 1
6 4 1 1 8 8 2 3 5 1 2
8 10 3 10 5 2 5 6 3 2 6 1
6 2 2 10 2 5 5 3 3 2 1
8 7 7 4 2 7 5 8 7 3 8 7 3
7 4 7 4 5 2 3 8 2 3 3
8 6 4 3 8 3 6 4 8 1 2
8 6 3 2 6 3 5 6 2 1 7
6 6 5 6 1 4 5 3 2 2 5
3 8 4 2 3 3 2 3
2 1 7 3 1 2 2
8 3 3 1 2 3 2 7 10 1 5
2 10 10 1 1
5 7 3 9 8 5 2 2 4
5 3 3 2 4 9 2 4 1
6 9 7 5 8 1 3 2 2 4
4 4 2 3 8 2 1 3
4 2 6 3 4 1 3
8 6 10 8 9 5 5 1 10 3 3 5 4
2 5 8 1 1
1 8 3 1 1 1
6 4 6 8 3 9 1 3 4 2 1
2 9 1 1 1
1 6 2 1 1
1 4 3 1 1 1
6 7 4 7 1 10 4 3 1 4 6
1 7 1 1
2 10 3 2 2 2
8 6 9 7 8 6 5 9 5 2 6 7
6 7 5 1 4 8 5 2 4 6
3 3 8 4 1 3
4 2 5 9 9 1 2
8 2 8 4 7 3 2 1 9 1 2
2 10 3 1 1
6 6 2 2 7 8 3 2 5 5
4 10 4 10 10 2 3 1
3 7 2 4 3 1 1 3
5 5 6 3 8 6 3 1 1 4
1 1 1 1
7 4 6 7 10 2 6 2 2 5 4
5 4 7 1 6 6 1 1
5 8 1 10 3 3 3 5 5 3
8 2 7 5 8 5 1 7 1 2 2 3
5 6 5 10 3 8 2 1 5
6 3 7 1 5 3 10 1 1
7 9 8 9 3 7 6 3 3 5 3 5
`

type testCase struct {
	n int
	a []int64
	m int
	k []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesDData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		pos := 0
		nextInt := func() (int, error) {
			if pos >= len(fields) {
				return 0, fmt.Errorf("line %d: unexpected end of line", idx+1)
			}
			v, err := strconv.Atoi(fields[pos])
			pos++
			if err != nil {
				return 0, fmt.Errorf("line %d: parse int: %w", idx+1, err)
			}
			return v, nil
		}

		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		if len(fields) < 1+n+1 {
			return nil, fmt.Errorf("line %d: insufficient data", idx+1)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("line %d: a[%d]: %w", idx+1, i, err)
			}
			a[i] = int64(val)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: m: %w", idx+1, err)
		}
		if len(fields) != pos+m {
			return nil, fmt.Errorf("line %d: expected %d ks, got %d", idx+1, m, len(fields)-pos)
		}
		k := make([]int, m)
		for i := 0; i < m; i++ {
			val, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("line %d: k[%d]: %w", idx+1, i, err)
			}
			k[i] = val
		}
		cases = append(cases, testCase{n: n, a: a, m: m, k: k})
	}
	return cases, nil
}

// solveCase mirrors the reference solution logic from 211D.go and returns answers for the given ks.
func solveCase(n int, a []int64, ks []int) []float64 {
	// compute L: distance to previous strictly less
	L := make([]int, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			L[i] = i + 1
		} else {
			L[i] = i - stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	// compute R: distance to next less or equal
	R := make([]int, n)
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			R[i] = n - i
		} else {
			R[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	// difference arrays for alpha (da) and beta (db)
	da := make([]int64, n+3)
	db := make([]int64, n+3)
	for i := 0; i < n; i++ {
		A := L[i]
		B := R[i]
		mn := A
		mx := B
		if B < A {
			mn = B
			mx = A
		}
		total := int64(A + B - 1)
		ai := a[i]
		if mn > 0 {
			da[1] += ai
			da[mn+1] -= ai
		}
		if mx > mn {
			v := ai * int64(mn)
			db[mn+1] += v
			db[mx+1] -= v
		}
		if int(mx)+1 <= int(total) {
			l := mx + 1
			r := int(total)
			da[l] -= ai
			da[r+1] += ai
			db[l] += ai * int64(A+B)
			db[r+1] -= ai * int64(A+B)
		}
	}
	// build prefix sums and compute S[k]
	S := make([]float64, n+1)
	var curA, curB int64
	for k := 1; k <= n; k++ {
		curA += da[k]
		curB += db[k]
		S[k] = float64(curA)*float64(k) + float64(curB)
	}
	ans := make([]float64, len(ks))
	for i, k := range ks {
		ans[i] = S[k] / float64(n-k+1)
	}
	return ans
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", tc.m)
	for _, v := range tc.k {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expect := solveCase(tc.n, tc.a, tc.k)
		input := buildInput(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotStr)
		if len(gotFields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers, got %d\n", idx+1, len(expect), len(gotFields))
			os.Exit(1)
		}
		for i, g := range gotFields {
			var v float64
			if _, err := fmt.Sscan(g, &v); err != nil {
				fmt.Fprintf(os.Stderr, "case %d parse output: %v\n", idx+1, err)
				os.Exit(1)
			}
			if fmt.Sprintf("%.10f", v) != fmt.Sprintf("%.10f", expect[i]) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f got %.10f (index %d)\n", idx+1, expect[i], v, i)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

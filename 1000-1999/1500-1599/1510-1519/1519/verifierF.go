package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	m int
	a []int
	b []int
	c [][]int
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
100
3 2
3 2 3
3 3
2 0
1 0
2 0
1 1
2
2
0
2 3
1 3
1 1 3
0 1 1
0 1 0
1 1
3
3
1
1 1
1
1
0
1 1
1
2
1
1 3
3
3 1 1
2 0 1
2 1
2 2
1
0
1
1 2
2
3 3
0 2
3 3
2 1 2
2 2 2
2 1 0
1 1 2
0 0 1
1 3
2
2 1 3
1 1 1
3 1
2 1 3
1
2
0
0
1 2
2
3 2
2 1
2 1
3 3
2
1
0
2 1
1 2
3
2
1
1 2
2
3 3
0 1
3 2
2 1 1
3 1
2 2
1 1
0 2
1 1
3
3
1
1 1
3
3
1
2 2
3 2
1 1
2 0
1 1
1 1
3
3
0
3 2
1 1 2
2 1
0 1
1 0
1 2
1 3
2
2 3 3
1 1 1
2 2
2 1
1 2
2 0
2 1
2 1
1 2
2
2
2
3 2
1 2 3
3 3
0 1
1 2
2 2
2 2
2 3
3 2
1 2
0 2
1 2
3
2 2
2 1
2 2
3 3
3 2
1 1
2 1
3 2
2 1 3
2 2
1 2
0 2
0 0
2 2
2 3
1 3
2 2
1 0
3 3
3 3 2
2 3 3
2 1 2
0 0 0
2 1 2
3 1
1 3 1
3
2
0
1
1 1
1
1
0
2 1
2 1
3
0
1
2 2
2 3
3 1
2 0
0 2
2 1
2 2
2
1
0
3 3
2 1 3
3 2 1
1 2 2
1 2 1
1 0 1
2 3
3 3
2 2 1
0 2 1
2 1 0
1 3
1
2 3 1
1 0 1
3 3
3 1 2
3 3 2
1 1 0
1 1 1
2 1 1
1 3
3
2 3 2
1 1 1
1 2
2
2 2
0 1
1 1
3
3
1
1 2
1
1 2
0 1
3 3
2 1 2
1 1 2
1 0 2
1 1 1
2 1 2
3 3
3 1 1
1 2 2
2 2 1
2 1 2
2 0 1
2 2
3 3
3 1
0 0
1 2
2 3
1 3
3 1 1
1 1 2
2 2 2
3 3
2 2 2
2 3 3
0 1 1
0 1 0
2 1 2
2 3
2 2
1 1 2
2 2 1
1 2 0
2 3
3 1
1 3 1
1 0 0
2 0 1
3 2
1 2 3
2 1
0 0
0 2
2 0
3 1
2 3 2
1
0
2
0
2 1
1 3
1
0
2
1 3
2
3 1 2
0 2 2
1 3
2
3 1 1
1 0 1
3 3
2 2 2
3 2 1
0 2 1
0 2 2
2 2 1
3 3
3 1 2
1 3 1
1 1 2
0 0 1
1 2 1
1 1
3
1
2
2 1
3 3
1
0
0
3 3
2 1 2
2 2 2
0 2 0
2 2 1
0 0 1
3 1
3 1 1
3
1
0
0
3 3
2 2 1
2 2 3
2 1 1
1 0 2
1 2 1
3 2
2 3 1
3 3
0 2
0 1
0 0
1 1
1
1
0
1 2
2
3 2
0 1
2 2
3 2
2 2
1 0
0 0
3 1
1 2 2
2
0
2
1
1 2
1
2 3
0 0
2 1
2 2
2
0
2
1 3
2
2 1 1
2 1 1
2 3
1 3
1 1 1
2 0 2
0 2 0
3 2
3 2 3
2 1
1 1
1 0
2 2
2 1
1 1
3
0
0
3 3
2 3 2
1 2 1
1 2 1
0 0 0
2 2 2
2 1
3 1
1
1
0
2 1
3 1
3
0
1
2 1
1 2
2
1
1
3 1
1 1 1
3
2
2
1
1 3
1
1 2 1
0 2 0
3 1
3 1 3
1
1
2
1
2 2
3 2
3 3
0 2
0 1
3 2
1 3 1
1 2
1 0
1 1
1 0
1 3
1
2 1 3
0 1 2
3 1
2 1 2
3
0
0
2
1 2
2
1 1
1 1
2 3
1 2
3 1 3
0 1 0
2 2 0
3 1
3 2 2
1
1
2
0
2 1
2 1
1
1
1
3 2
2 1 1
3 3
2 0
0 0
0 0
2 2
2 3
1 1
2 1
0 1
1 2
2
1 1
2 2
1 3
1
2 3 3
0 2 2
1 2
2
2 1
1 1
2 1
2 3
2
0
1
1 3
1
3 2 2
2 0 0
1 2
3
1 1
0 0
3 1
2 3 2
2
1
0
0
2 1
2 3
3
1
0
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
		a := make([]int, n)
		b := make([]int, m)
		for j := 0; j < n; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad a[%d]: %v", i+1, j, err)
			}
			a[j] = v
		}
		for j := 0; j < m; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad b[%d]: %v", i+1, j, err)
			}
			b[j] = v
		}
		c := make([][]int, n)
		for r := 0; r < n; r++ {
			c[r] = make([]int, m)
			for col := 0; col < m; col++ {
				v, err := nextInt()
				if err != nil {
					return nil, fmt.Errorf("case %d bad c[%d][%d]: %v", i+1, r, col, err)
				}
				c[r][col] = v
			}
		}
		res = append(res, testCase{n: n, m: m, a: a, b: b, c: c})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of data")
	}
	return res, nil
}

// solve mirrors 1519F.go.
func solve(tc testCase) int64 {
	n, m := tc.n, tc.m
	a := tc.a
	b := tc.b
	c := tc.c

	sumA, sumB := 0, 0
	for _, v := range a {
		sumA += v
	}
	for _, v := range b {
		sumB += v
	}
	if sumA > sumB {
		return -1
	}

	pow := make([]int, n)
	pow[0] = 1
	for i := 1; i < n; i++ {
		pow[i] = pow[i-1] * 5
	}
	finalNeed := 0
	for i := 0; i < n; i++ {
		finalNeed += a[i] * pow[i]
	}

	type state struct {
		need int
		v1   int
		v2   int
		rem  int
	}

	start := state{0, 0, 0, 0}
	dp := map[state]int{start: 0}
	best := int64(1<<63 - 1)

	for len(dp) > 0 {
		next := make(map[state]int)
		for st, cost := range dp {
			if int64(cost) >= best {
				continue
			}
			if st.v2 == m {
				if st.need == finalNeed && int64(cost) < best {
					best = int64(cost)
				}
				continue
			}
			i := st.v1
			j := st.v2
			cur := (st.need / pow[i]) % 5
			maxF := a[i] - cur
			if maxF > b[j]-st.rem {
				maxF = b[j] - st.rem
			}
			for f := 0; f <= maxF; f++ {
				nn := st.need + f*pow[i]
				nc := cost
				if f > 0 {
					nc += c[i][j]
				}
				nv1 := i + 1
				nv2 := j
				nr := st.rem + f
				if nv1 == n {
					nv1 = 0
					nv2 = j + 1
					nr = 0
				}
				nst := state{nn, nv1, nv2, nr}
				if nv2 == m {
					if nn == finalNeed && int64(nc) < best {
						best = int64(nc)
					}
					continue
				}
				if prev, ok := next[nst]; !ok || nc < prev {
					next[nst] = nc
				}
			}
		}
		dp = next
	}

	if best == int64(1<<63-1) {
		return -1
	}
	return best
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.c[i][j]))
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strconv.FormatInt(expect, 10) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

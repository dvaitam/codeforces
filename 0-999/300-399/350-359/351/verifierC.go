package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (previously in testcasesC.txt) so the verifier is self contained.
const rawTestcasesC = `
1 4 2 6
2 10 5 10 4 10
1 6 7 7
5 12 9 8 9 5 1 1 6 8 6 7
4 6 9 3 4 4 1 3 6 3
2 12 9 9 3 8
4 12 10 6 6 8 3 7 8 9
2 16 5 8 9 9
3 16 8 6 10 9 8 8
2 12 3 10 5 8
3 10 9 9 9 9 10 10
4 10 4 8 9 6 10 2 6 1
2 4 1 10 1 5
5 8 2 9 3 5 4 4 1 7 1 1
3 12 3 4 1 2 2 2
1 2 1 6
3 6 3 3 9 1 7 10
1 8 3 1
1 12 10 2
3 12 8 1 5 8 9 10
1 10 7 10
2 16 4 2 6 2
1 16 3 9
5 14 8 9 6 3 6 5 5 10 7 1
5 6 1 5 1 3 3 3 2 8 4 9
1 8 4 8
1 10 2 10
2 12 5 7 5 9
1 6 1 7
4 6 2 9 2 4 2 2 1 3
2 4 4 1 9 8
4 10 9 7 4 4 7 7 9 1
5 2 7 9 10 3 2 8 6 1 9 2
5 12 5 6 5 1 7 2 2 5 4 1
4 2 7 8 8 4 10 10 2 1
3 2 6 5 2 4 8 4
1 12 7 8
2 12 7 2 5 2
1 4 10 6
4 8 2 1 10 8 1 8 5 6
4 6 6 5 8 9 8 7 8 5
4 8 3 8 10 5 9 7 2 10
5 4 2 6 3 9 3 7 2 2 1 3
3 14 4 6 8 3 9 5
1 6 9 7
1 12 9 4
5 10 3 3 8 4 7 6 10 3 8 8
1 14 3 7
5 2 8 5 7 5 7 8 6 9 6 2
2 8 7 7 1 6
4 16 3 2 1 7 4 10 10 7
2 4 7 9 4 5
5 8 8 10 3 1 10 7 8 5 9 10
2 16 4 2 6 1
4 4 10 8 6 8 5 9 8 1
1 12 3 7
3 6 1 3 8 7 8 5
2 2 5 9 8 1
3 2 9 7 10 8 4 5
4 6 8 9 5 2 5 6 5 6
3 14 9 2 9 4 7 10
5 6 9 2 5 1 4 8 9 4 9 5
1 4 2 7
3 8 6 6 2 6 8 6
2 16 8 5 8 3
4 8 5 6 3 2 4 8 4 6
2 12 3 3 4 5
5 14 7 6 5 10 9 10 6 7 5 9
5 4 6 5 7 8 3 5 6 8 8 2
2 12 7 3 1 2
3 6 6 2 7 1 9 6
2 14 9 5 8 3
3 12 4 8 2 3 4 6
3 6 7 6 5 2 6 4
2 8 10 1 6 6
5 2 3 3 2 7 8 5 3 6 9 10
1 12 10 7
2 2 7 8 8 10
3 4 10 9 9 8 7 8
2 14 7 9 8 1
1 16 10 3
1 6 2 7
3 16 1 5 2 6 4 3
1 6 7 2
3 16 1 8 4 2 8 3
5 2 3 9 9 1 1 4 9 1 9 6
5 8 3 6 8 1 3 9 2 4 2 8
2 2 10 4 7 6
5 14 9 9 3 9 2 3 4 3 7 4
3 12 7 3 7 3 7 6
3 4 9 2 8 5 5 9
4 10 4 7 3 9 2 1 10 9
2 8 4 7 10 1
2 2 5 8 9 1
2 6 10 6 1 4
1 6 9 3
1 16 5 4
2 12 5 9 10 2
4 14 1 8 5 2 5 1 4 7
3 10 9 7 10 9 4 7
`

type testCase struct {
	n int
	k int64
	a []int64
	b []int64
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcasesC, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: expected at least n and k", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		k, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %w", idx+1, err)
		}
		if len(fields) != 2+2*n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 2*n, len(fields)-2)
		}
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %w", idx+1, i, err)
			}
			a[i] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[2+n+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b[%d]: %w", idx+1, i, err)
			}
			b[i] = v
		}
		cases = append(cases, testCase{n: n, k: k, a: a, b: b})
	}
	return cases, nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// mult returns A*B under min-plus.
func mult(A, B [][]int64, INF int64) [][]int64 {
	n := len(A)
	C := make([][]int64, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			C[i][j] = INF
			for k := 0; k < n; k++ {
				v := A[i][k] + B[k][j]
				if v < C[i][j] {
					C[i][j] = v
				}
			}
		}
	}
	return C
}

// solve351C mirrors the logic in 351C.go.
func solve351C(tc testCase) int64 {
	n, m := tc.n, tc.k
	a := tc.a
	b := tc.b
	N := n + 1
	INF := int64(9e18)
	C := make([][]int64, N)
	for i := 0; i < N; i++ {
		C[i] = make([]int64, N)
		for j := 0; j < N; j++ {
			C[i][j] = INF
		}
	}
	for u := 0; u < n; u++ {
		maxh := n - u
		dp := make([]int64, maxh+1)
		ndp := make([]int64, maxh+1)
		for h := 0; h <= maxh; h++ {
			dp[h] = INF
		}
		dp[0] = 0
		for i := 0; i < n; i++ {
			for h := 0; h <= maxh; h++ {
				ndp[h] = INF
			}
			for h := 0; h <= maxh; h++ {
				cost := dp[h]
				if cost >= INF {
					continue
				}
				h2 := h + 1
				if h2 > maxh {
					h2 = maxh
				}
				ndp[h2] = min(ndp[h2], cost+a[i])
				if h > 0 {
					h3 := h - 1
					ndp[h3] = min(ndp[h3], cost+b[i])
				}
			}
			dp, ndp = ndp, dp
		}
		for h := 0; h <= maxh; h++ {
			cost := dp[h]
			if cost >= INF {
				continue
			}
			f := u + h
			v := f
			if f >= n {
				v = n
			}
			if C[u][v] > cost {
				C[u][v] = cost
			}
		}
	}
	base := n
	sizeH := 2*n + 1
	off := n
	dp2 := make([]int64, sizeH)
	ndp2 := make([]int64, sizeH)
	for i := 0; i < sizeH; i++ {
		dp2[i] = INF
	}
	dp2[off] = 0
	for i := 0; i < n; i++ {
		for j := 0; j < sizeH; j++ {
			ndp2[j] = INF
		}
		for hi := 0; hi < sizeH; hi++ {
			cost := dp2[hi]
			if cost >= INF {
				continue
			}
			h := hi - off
			h2 := h + 1
			if h2 > n {
				h2 = n
			}
			hi2 := h2 + off
			ndp2[hi2] = min(ndp2[hi2], cost+a[i])
			h3 := h - 1
			if h3 < -n {
				h3 = -n
			}
			hi3 := h3 + off
			ndp2[hi3] = min(ndp2[hi3], cost+b[i])
		}
		dp2, ndp2 = ndp2, dp2
	}
	for hi := 0; hi < sizeH; hi++ {
		cost := dp2[hi]
		if cost >= INF {
			continue
		}
		h := hi - off
		f := base + h
		v := f
		if f < n {
			v = f
		} else {
			v = n
		}
		if C[n][v] > cost {
			C[n][v] = cost
		}
	}

	res := make([][]int64, N)
	for i := 0; i < N; i++ {
		res[i] = make([]int64, N)
		for j := 0; j < N; j++ {
			if i == j {
				res[i][j] = 0
			} else {
				res[i][j] = INF
			}
		}
	}
	baseM := C
	for m > 0 {
		if m&1 == 1 {
			res = mult(res, baseM, INF)
		}
		baseM = mult(baseM, baseM, INF)
		m >>= 1
	}
	return res[0][0]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expect := solve351C(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: failed to parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

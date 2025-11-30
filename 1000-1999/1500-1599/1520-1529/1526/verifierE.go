package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

type testCase struct {
	n   int
	k   int
	perm []int
}

// Embedded testcases from testcasesE.txt.
const testcaseData = `
100
1 2
0
10 7
1 2 0 5 6 4 8 9 7 3
1 9
0
8 2
4 3 0 1 2 7 5 6
8 10
2 1 4 0 6 7 5 3
5 9
1 0 4 3 2
8 2
4 0 1 6 3 5 2 7
5 4
0 2 1 4 3
3 8
2 0 1
5 10
0 3 2 4 1
2 3
1 0
9 10
1 5 0 2 3 4 6 7 8
7 6
1 6 5 4 0 3 2
7 2
2 4 6 1 5 3 0
1 3
0
8 10
4 5 3 1 2 6 7 0
4 4
3 0 2 1
1 8
0
8 6
1 7 3 4 6 0 5 2
7 5
1 4 5 2 3 0 6
2 6
0 1
9 4
7 5 0 1 2 6 4 3 8
6 2
5 4 1 3 2 0
9 3
2 6 4 7 8 0 3 5 1
10 5
1 2 7 4 8 0 6 3 5 9
6 2
4 3 1 2 0 5
1 2
0
7 4
1 5 3 0 4 6 2
6 6
4 1 3 0 5 2
5 6
4 1 3 0 2
6 4
2 4 5 3 0 1
3 3
1 0 2
9 2
8 7 3 1 2 5 0 6 4
8 9
6 0 5 4 2 7 3 1
7 2
1 4 0 6 5 3 2
6 7
1 0 5 2 4 3
4 2
1 2 0 3
4 1
3 0 1 2
2 10
1 0
6 3
0 4 5 3 2 1
5 5
1 0 4 3 2
8 3
5 0 6 2 4 3 1 7
1 3
0
6 6
3 4 2 0 1 5
6 10
1 4 2 3 5 0
1 8
0
8 2
6 2 7 4 5 1 0 3
1 5
0
5 1
4 3 1 0 2
7 2
0 4 2 1 3 6 5
8 9
3 7 1 4 6 2 0 5
1 4
0
9 1
7 4 2 5 8 3 0 6 1
10 7
2 7 3 6 5 1 4 9 8 0
8 4
5 2 1 3 6 7 4 0
5 8
1 4 3 0 2
8 8
6 4 5 0 1 3 2 7
7 8
6 3 2 0 5 4 1
10 6
5 0 8 9 1 6 7 4 2 3
5 7
3 4 2 0 1
7 3
2 1 0 4 3 5 6
7 6
0 4 1 5 3 2 6
1 4
0
8 6
5 4 1 7 6 2 0 3
10 7
4 9 7 2 1 8 0 5 6 3
5 2
2 3 0 4 1
6 5
4 2 0 3 5 1
8 10
7 4 6 5 3 2 0 1
6 5
0 3 4 2 5 1
6 8
0 4 2 1 5 3
4 5
0 1 2 3
8 1
7 4 1 0 2 3 6 5
6 8
4 1 3 2 5 0
10 7
8 0 6 4 3 9 2 1 7 5
7 9
3 1 4 6 0 5 2
9 10
1 4 7 2 0 8 5 3 6
10 9
5 9 8 2 7 4 0 6 1 3
1 8
0
3 9
1 0 2
10 10
5 0 8 1 7 2 3 4 6 9
4 4
1 0 2 3
4 5
3 1 0 2
3 2
2 0 1
4 10
2 0 3 1
8 5
3 4 6 7 2 0 5 1
6 2
0 5 1 4 2 3
5 4
1 0 3 2 4
6 7
4 3 5 1 0 2
10 10
3 0 7 1 4 2 9 5 8 6
2 3
1 0
5 1
0 4 3 2 1
9 6
2 5 6 3 4 8 0 7 1
7 10
2 5 3 0 1 6 4
7 9
6 0 4 3 1 5 2
10 10
3 5 1 9 0 2 6 7 8 4
8 3
5 2 7 1 3 6 0 4
9 3
8 0 3 6 4 7 5 1 2
2 9
0 1
6 7
3 5 4 0 1 2
10 10
4 9 5 3 6 0 2 7 1 8
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		return nil, fmt.Errorf("no test data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	if len(lines)-1 != t*2 {
		return nil, fmt.Errorf("expected %d lines after t, got %d", t*2, len(lines)-1)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for i := 0; i < t; i++ {
		parts := strings.Fields(lines[idx])
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: bad n k line", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		permStr := strings.Fields(lines[idx+1])
		if len(permStr) != n {
			return nil, fmt.Errorf("case %d expected %d perm values, got %d", i+1, n, len(permStr))
		}
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(permStr[j])
			if err != nil {
				return nil, fmt.Errorf("case %d bad perm %d: %v", i+1, j+1, err)
			}
			perm[j] = v
		}
		res = append(res, testCase{n: n, k: k, perm: perm})
		idx += 2
	}
	return res, nil
}

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

// solve mirrors 1526E.go.
func solve(tc testCase) string {
	n, k := tc.n, tc.k
	p := tc.perm
	rank := make([]int, n+1)
	for i := 0; i < n; i++ {
		rank[p[i]] = i
	}
	rank[n] = -1

	strict := 0
	for idx := 0; idx < n-1; idx++ {
		i, j := p[idx], p[idx+1]
		ri1 := -1
		if i+1 < n {
			ri1 = rank[i+1]
		}
		rj1 := -1
		if j+1 < n {
			rj1 = rank[j+1]
		}
		if ri1 > rj1 {
			strict++
		}
	}

	if k-strict <= 0 {
		return "0"
	}

	N := int64(n + k - strict - 1)
	fact := make([]int64, N+1)
	invFact := make([]int64, N+1)
	fact[0] = 1
	for i := int64(1); i <= N; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[N] = modPow(fact[N], MOD-2)
	for i := N; i > 0; i-- {
		invFact[i-1] = invFact[i] * i % MOD
	}

	R := int64(n)
	ans := fact[N]
	ans = ans * invFact[R] % MOD
	ans = ans * invFact[N-R] % MOD
	return strconv.FormatInt(ans, 10)
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

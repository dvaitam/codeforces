package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 1_000_000_007

const testcasesB1 = `100
1 1
6 1
1 2
2 3
3 4
1 5
4 6
5 2
1 2
2 3
3 4
3 5
5 3
1 2
2 3
1 4
1 5
1 1
4 1
1 2
1 3
3 4
4 1
1 2
2 3
3 4
2 2
1 2
1 1
2 1
1 2
3 2
1 2
1 3
5 3
1 2
1 3
3 4
4 5
6 3
1 2
2 3
3 4
4 5
2 6
2 1
1 2
1 1
3 3
1 2
1 3
5 3
1 2
2 3
2 4
3 5
6 2
1 2
2 3
2 4
1 5
3 6
2 1
1 2
5 2
1 2
2 3
2 4
1 5
2 1
1 2
1 1
6 1
1 2
1 3
3 4
4 5
3 6
6 2
1 2
2 3
1 4
1 5
1 6
6 1
1 2
1 3
3 4
4 5
3 6
5 3
1 2
2 3
2 4
3 5
1 1
5 1
1 2
1 3
3 4
2 5
3 1
1 2
1 3
2 2
1 2
2 1
1 2
5 1
1 2
1 3
1 4
1 5
5 2
1 2
1 3
1 4
4 5
4 1
1 2
1 3
1 4
6 1
1 2
2 3
1 4
4 5
4 6
1 1
2 2
1 2
3 2
1 2
1 3
5 3
1 2
1 3
1 4
2 5
1 1
3 2
1 2
1 3
4 3
1 2
1 3
2 4
2 1
1 2
4 3
1 2
2 3
3 4
6 2
1 2
2 3
2 4
3 5
5 6
4 3
1 2
2 3
3 4
5 1
1 2
1 3
3 4
1 5
6 1
1 2
1 3
1 4
2 5
5 6
6 3
1 2
2 3
1 4
4 5
4 6
4 2
1 2
2 3
2 4
4 3
1 2
2 3
1 4
1 1
3 2
1 2
1 3
4 3
1 2
2 3
2 4
2 1
1 2
2 2
1 2
4 1
1 2
1 3
1 4
3 1
1 2
1 3
2 1
1 2
3 2
1 2
1 3
3 1
1 2
1 3
1 1
3 3
1 2
1 3
2 2
1 2
4 3
1 2
2 3
3 4
4 3
1 2
1 3
3 4
6 2
1 2
1 3
2 4
1 5
2 6
4 3
1 2
2 3
2 4
2 1
1 2
1 1
3 2
1 2
1 3
5 2
1 2
1 3
3 4
3 5
2 2
1 2
4 2
1 2
1 3
3 4
6 1
1 2
1 3
1 4
2 5
2 6
1 1
1 1
1 1
1 1
5 3
1 2
2 3
2 4
4 5
5 1
1 2
2 3
2 4
4 5
6 3
1 2
2 3
2 4
3 5
5 6
6 1
1 2
2 3
3 4
4 5
1 6
5 1
1 2
1 3
3 4
2 5
1 1
4 1
1 2
1 3
1 4
1 1
6 3
1 2
1 3
2 4
1 5
2 6
3 2
1 2
2 3
1 1
1 1
6 1
1 2
1 3
3 4
1 5
2 6
6 3
1 2
2 3
3 4
1 5
3 6
6 3
1 2
2 3
2 4
4 5
5 6
1 1
6 3
1 2
1 3
3 4
1 5
1 6
1 1
3 2
1 2
1 3
3 2
1 2
1 3
1 1`

type testCase struct {
	input    string
	expected string
}

func powmod(a, e int64) int64 {
	a %= MOD
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func expected(n, k int, edges [][2]int) string {
	if k == 1 || k == 3 {
		return "1"
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	sz := make([]int, n+1)
	var sumDist int64
	var dfs func(int, int)
	dfs = func(v, p int) {
		sz[v] = 1
		for _, to := range adj[v] {
			if to == p {
				continue
			}
			dfs(to, v)
			s := sz[to]
			sumDist += int64(s) * int64(n-s)
			sz[v] += s
		}
	}
	dfs(1, 0)
	numerator := (2 * (sumDist % MOD)) % MOD
	denom := int64(n) * int64(n-1) % MOD
	invDenom := powmod(denom, MOD-2)
	ans := (1 + numerator*invDenom%MOD) % MOD
	return fmt.Sprintf("%d", ans)
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcasesB1)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n/k", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		k, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad k: %w", caseIdx+1, err)
		}
		pos++
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing edge", caseIdx+1)
			}
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge u: %w", caseIdx+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge v: %w", caseIdx+1, err)
			}
			edges[j] = [2]int{u, v}
			pos += 2
		}
		var sb bytes.Buffer
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: expected(n, k, edges),
		})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB1 /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(tc.input)
		res, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, res)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(res))
		if got != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

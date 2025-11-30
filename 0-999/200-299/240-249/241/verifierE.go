package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt.
const testcasesEData = `
5 6 2 4 3 4 1 5 2 3 4 5 3 5
5 3 4 5 2 5 3 5
3 2 2 3 1 2
4 2 2 3 3 4
4 5 2 4 3 4 1 4 2 3 1 3
3 3 2 3 1 2 1 3
3 1 1 3
2 1 1 2
5 3 4 5 1 2 3 4
4 5 2 4 1 2 3 4 1 4 2 3
5 6 2 4 3 4 2 3 2 5 1 3 3 5
5 2 2 5 1 4
5 6 3 4 1 5 2 3 4 5 2 5 1 3
4 2 2 4 3 4
5 2 3 4 1 4
5 3 4 5 2 4 3 5
3 2 1 2 1 3
3 2 2 3 1 2
4 4 2 4 1 2 1 3 3 4
5 1 4 5
2 1 1 2
2 1 1 2
3 3 2 3 1 2 1 3
4 2 1 3 3 4
3 3 2 3 1 2 1 3
4 4 2 3 2 4 1 2 3 4
2 1 1 2
4 6 2 4 1 2 3 4 1 4 2 3 1 3
3 2 2 3 1 3
5 5 1 2 1 5 1 4 4 5 3 5
5 2 2 5 1 3
5 6 1 4 2 3 4 5 2 5 1 3 3 5
4 6 2 4 1 2 3 4 1 4 2 3 1 3
4 4 2 4 1 3 3 4 1 4
2 1 1 2
2 1 1 2
4 5 2 4 1 2 3 4 2 3 1 3
5 3 4 5 2 5 1 3
4 3 2 3 2 4 3 4
4 4 2 3 2 4 1 2 1 4
4 2 2 3 2 4
4 6 2 4 3 4 1 2 3 4 2 4 1 3
4 5 3 4 2 3 1 5 3 4
2 1 1 2
2 1 1 2
2 1 1 2
2 1 1 2
3 2 2 3 2 3
2 1 1 2
4 2 2 4 1 2
2 1 1 2
4 2 2 3 2 4
3 3 1 2 2 3
2 1 1 2
4 4 2 4 2 3 2 4 1 4
2 1 1 2
2 1 1 2
2 1 1 2
4 2 1 3 1 2
2 1 1 2
2 1 1 2
2 1 1 2
2 1 1 2
2 1 1 2
2 1 1 2
2 1 1 2
3 2 1 2 2 3
2 1 1 2
3 3 2 3 3 1 1 3
3 3 3 1 2 3 1 2
4 3 1 3 2 4 1 3
2 1 1 2
2 1 1 2
2 1 1 2
4 4 1 3 3 4 2 3 2 4
5 5 1 3 2 3 3 4 2 5 4 5
3 2 1 2 3 2
3 2 1 2 1 3
3 2 1 3 1 2
4 4 2 4 1 2 2 3 3 4
5 6 2 4 3 4 1 5 1 2 3 4 3 4
3 3 3 1 2 3 2 3
4 3 2 4 2 3 2 4
4 6 3 5 1 2 2 4 2 6 3 6
2 1 1 2
2 1 1 2
4 3 2 4 1 2 3 4
4 3 2 3 1 4 3 4
3 3 2 3 2 3 1 2
3 2 2 3 1 3
4 5 1 4 2 3 2 5 1 3 3 4
5 4 1 4 2 4 2 3 2 5
2 1 1 2
4 3 2 3 1 4 1 2
3 3 2 3 2 3 1 3
5 4 2 4 1 4 3 4 2 4
3 2 1 2 1 3
2 1 1 2
3 2 1 2 1 3
`

type testCase struct {
	n int
	m int
	u []int
	v []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesEData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		if len(fields) != 2+2*m {
			return nil, fmt.Errorf("line %d: expected %d edge values, got %d", idx+1, 2*m, len(fields)-2)
		}
		u := make([]int, m)
		v := make([]int, m)
		pos := 2
		for i := 0; i < m; i++ {
			ui, err := strconv.Atoi(fields[pos+i*2])
			if err != nil {
				return nil, fmt.Errorf("line %d: u[%d]: %w", idx+1, i, err)
			}
			vi, err := strconv.Atoi(fields[pos+i*2+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: v[%d]: %w", idx+1, i, err)
			}
			u[i] = ui - 1
			v[i] = vi - 1
		}
		cases = append(cases, testCase{n: n, m: m, u: u, v: v})
	}
	return cases, nil
}

// solve replicates the logic from 241E.go for a single test case.
func solve(tc testCase) string {
	n, m := tc.n, tc.m
	u, v := tc.u, tc.v
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		adj[u[i]] = append(adj[u[i]], v[i])
	}
	vis := make([]bool, n)
	reach := make([]bool, n)
	var dfs func(int)
	dfs = func(x int) {
		vis[x] = true
		if x == n-1 {
			reach[x] = true
		}
		for _, w := range adj[x] {
			if !vis[w] {
				dfs(w)
			}
			if reach[w] {
				reach[x] = true
			}
		}
	}
	dfs(0)

	s := make([]int, n)
	updated := true
	for l := 1; l <= n && updated; l++ {
		updated = false
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if reach[ui] && reach[vi] {
				diff := s[ui] - s[vi]
				if diff < 1 {
					s[vi] = s[ui] - 1
					updated = true
				} else if diff > 2 {
					s[ui] = s[vi] + 2
					updated = true
				}
			}
		}
	}
	var out strings.Builder
	if updated {
		out.WriteString("No")
	} else {
		out.WriteString("Yes\n")
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if reach[ui] && reach[vi] {
				out.WriteString(fmt.Sprintf("%d\n", s[ui]-s[vi]))
			} else {
				out.WriteString("1\n")
			}
		}
		if out.Len() > 0 && out.String()[out.Len()-1] == '\n' {
			// leave trailing newline for comparison trimming.
		}
	}
	return strings.TrimSpace(out.String())
}

func runCandidate(bin string, tc testCase) (string, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.m; i++ {
		fmt.Fprintf(&input, "%d %d\n", tc.u[i]+1, tc.v[i]+1)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		exp := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

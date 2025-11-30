package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `100
8 6
4 8
2 7
2 5
1 5
7 8
2 3
7 6
1 7
1 6
3 4
4 7
2 4
2 7
3 3
1 2
1 3
2 3
8 6
5 6
3 5
2 3
3 8
1 4
3 7
6 6
3 5
2 6
5 6
1 2
1 4
2 3
6 7
4 6
3 6
2 5
1 5
1 2
2 3
2 6
6 6
3 5
1 5
5 6
2 4
1 4
2 3
6 5
1 4
1 6
3 6
2 5
2 6
6 10
3 5
2 4
1 5
1 4
4 5
1 6
4 6
3 4
2 5
3 6
8 4
1 5
3 8
5 6
2 8
5 1
3 5
4 5
1 4
1 3
1 2
2 4
2 3
8 2
5 8
1 3
3 3
1 3
2 3
1 2
7 17
3 5
3 6
3 4
1 4
1 2
1 6
2 4
4 7
1 5
2 5
2 7
1 7
5 7
4 5
5 6
2 3
3 7
8 3
7 8
3 4
1 8
3 2
1 2
2 3
6 15
2 3
1 4
3 4
1 6
1 5
1 2
5 6
1 3
4 5
2 5
3 6
2 6
2 4
4 6
3 5
8 3
3 4
6 8
1 7
3 3
1 3
1 2
2 3
6 10
3 4
4 6
2 3
5 6
1 4
1 6
1 3
2 6
2 4
3 5
6 3
2 5
4 5
1 5
8 28
1 2
2 3
1 5
3 6
3 8
4 5
3 4
2 7
2 6
1 6
4 7
5 7
2 4
1 8
2 5
3 7
6 7
1 7
4 6
3 5
2 8
1 3
4 8
6 8
7 8
1 4
5 6
5 8
3 0
4 1
2 4
3 3
1 3
1 2
2 3
7 12
6 7
4 7
1 2
2 6
2 5
1 3
1 7
4 5
1 4
1 5
2 4
2 7
7 4
2 6
3 4
3 6
3 5
7 11
4 6
3 7
3 5
1 5
2 3
3 4
4 7
5 6
1 7
1 4
3 6
4 6
1 4
3 4
1 3
2 4
1 2
2 3
3 2
1 2
1 3
8 16
1 7
1 4
1 2
4 5
3 7
3 6
7 8
1 3
3 5
2 5
6 8
3 8
2 7
6 7
5 6
4 6
7 13
3 7
2 6
3 5
1 2
1 5
2 4
5 6
4 7
2 5
3 4
1 4
2 3
4 6
7 13
5 6
1 3
3 5
1 7
3 6
1 2
3 7
2 6
4 5
1 6
2 7
1 5
2 3
7 1
4 7
8 15
6 8
3 8
1 2
4 5
5 8
1 3
2 7
3 4
2 6
1 4
4 6
4 7
2 3
1 6
5 7
6 7
4 5
3 6
1 4
4 6
5 6
2 5
3 5
3 1
1 2
8 14
1 3
5 8
7 8
2 5
1 5
5 7
2 8
4 7
2 7
6 8
6 7
4 5
1 2
3 6
6 4
1 6
5 6
1 4
1 3
6 14
1 6
2 5
4 6
2 3
1 3
3 6
5 6
1 5
1 4
3 4
1 2
2 6
2 4
3 5
7 11
1 7
1 6
3 7
4 6
2 4
1 2
5 7
4 7
2 6
3 6
3 5
7 11
3 7
3 6
2 3
1 4
4 5
2 4
4 7
1 7
5 7
1 2
3 4
7 20
1 4
3 7
1 3
4 5
6 7
5 7
2 5
1 7
2 7
4 7
2 6
2 3
3 4
1 5
3 5
3 6
1 2
5 6
4 6
2 4
3 3
2 3
1 3
1 2
8 2
2 7
2 5
3 2
1 2
1 3
4 5
3 4
1 4
2 3
1 2
1 3
7 21
1 6
2 6
2 4
1 4
6 7
3 4
3 5
3 7
1 7
1 5
3 6
5 7
1 2
4 5
2 3
1 3
2 7
4 7
5 6
2 5
4 6
4 2
2 4
1 3
5 8
3 4
3 5
2 4
2 3
4 5
1 2
1 5
1 3
5 1
1 4
5 0
8 9
5 8
3 7
2 3
2 4
3 8
2 6
1 6
1 4
1 8
3 0
3 0
5 8
2 3
3 4
1 3
1 4
1 2
2 4
4 5
1 5
5 8
4 5
3 5
1 2
3 4
1 5
1 4
2 3
2 4
3 2
1 2
2 3
3 2
1 2
2 3
5 5
1 4
2 4
1 2
4 5
3 4
8 27
5 6
3 8
2 8
3 7
1 8
5 7
5 8
1 6
1 2
6 8
1 7
2 7
2 3
2 6
4 5
7 8
1 4
4 7
2 5
4 6
6 7
1 3
1 5
4 8
2 4
3 4
3 6
7 9
1 6
2 4
1 5
4 5
1 7
1 3
3 6
3 5
3 7
8 10
1 7
2 3
1 6
5 7
5 6
1 8
4 8
1 3
3 6
1 4
8 12
7 8
1 2
1 6
2 6
3 7
1 3
4 5
6 7
5 6
2 3
1 4
3 6
5 8
4 5
1 5
3 5
2 5
1 3
1 4
2 3
3 4
7 15
1 5
3 7
2 4
5 7
2 6
5 6
2 7
2 5
1 2
1 4
1 6
6 7
4 7
3 6
3 4
3 0
6 15
2 4
1 3
5 6
2 3
2 6
1 4
1 6
1 5
1 2
3 6
3 5
4 5
2 5
3 4
4 6
8 3
3 4
3 7
2 8
3 1
1 2
6 12
3 6
2 4
4 6
4 5
1 6
1 2
3 4
3 5
5 6
2 6
1 3
1 5
8 26
6 7
5 7
3 4
3 6
3 7
7 8
1 5
5 6
3 5
3 8
2 7
2 3
1 7
4 6
6 8
1 4
1 8
1 3
2 4
2 5
1 6
5 8
2 8
4 7
4 8
1 2
7 12
1 4
3 4
1 6
4 6
3 5
2 5
5 7
4 7
6 7
1 7
3 6
2 3
8 22
2 7
3 7
2 5
2 8
5 6
6 8
6 7
3 4
5 8
2 3
3 8
3 5
4 8
1 4
1 8
4 6
2 4
3 6
1 3
5 7
4 7
7 8
5 4
1 4
3 5
1 2
2 5
6 13
1 3
1 4
1 6
3 6
4 5
2 6
1 5
4 6
1 2
5 6
2 4
3 4
2 5
7 17
2 7
1 3
4 7
3 5
1 5
1 2
6 7
3 7
5 6
4 6
2 6
1 7
2 4
5 7
2 5
1 4
4 5
4 6
3 4
1 4
2 4
1 3
1 2
2 3
4 5
1 2
1 3
3 4
2 4
1 4
8 18
3 7
5 6
6 8
1 3
4 7
2 5
1 6
2 3
2 4
1 5
3 4
1 4
2 7
3 5
1 8
4 6
2 8
4 8
8 20
2 6
2 7
1 6
5 8
3 7
6 8
3 5
4 5
1 4
1 8
4 6
3 6
3 8
1 5
2 4
5 6
2 3
1 7
4 8
1 2
3 3
1 2
2 3
1 3
4 5
1 4
1 2
2 4
1 3
2 3
5 10
4 5
1 3
1 2
2 5
2 4
3 4
1 5
1 4
3 5
2 3
5 2
2 4
2 3
5 6
3 5
1 5
1 3
2 4
2 5
4 5
7 12
4 7
1 3
1 2
3 6
1 4
5 6
1 7
1 6
3 4
2 5
2 3
4 5
4 6
1 4
2 4
1 2
1 3
3 4
2 3
6 7
2 3
4 5
4 6
1 3
3 5
3 6
1 6
4 1
2 4
8 24
2 4
6 8
3 7
2 6
1 6
4 7
3 4
6 7
1 4
3 8
2 8
1 3
4 6
1 5
5 8
5 7
1 7
2 3
2 7
3 6
4 5
3 5
7 8
1 8
4 0
8 22
3 8
1 3
1 4
2 3
3 7
6 8
3 5
5 8
4 6
3 4
6 7
3 6
1 8
2 7
5 6
1 7
2 6
2 8
4 7
4 5
1 6
7 8
8 2
2 5
5 8
3 0
4 0
5 9
2 4
3 4
4 5
2 5
2 3
1 4
1 5
1 2
3 5
5 1
1 3
7 4
3 5
1 2
2 6
4 7`

func solveCase(n, m int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	vis := make([]bool, n+1)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if vis[i] {
			continue
		}
		q = append(q[:0], i)
		vis[i] = true
		nodes := 0
		edgesCnt := 0
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			nodes++
			edgesCnt += len(adj[v])
			for _, u := range adj[v] {
				if !vis[u] {
					vis[u] = true
					q = append(q, u)
				}
			}
		}
		if edgesCnt/2 != nodes*(nodes-1)/2 {
			return "NO"
		}
	}
	return "YES"
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

type testCase struct {
	n     int
	m     int
	edges [][2]int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Fields(strings.TrimSpace(testcasesRaw))
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	pos := 0
	t, err := strconv.Atoi(lines[pos])
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+2 > len(lines) {
			return nil, fmt.Errorf("case %d: missing n m", caseIdx+1)
		}
		n, err := strconv.Atoi(lines[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %w", caseIdx+1, err)
		}
		m, err := strconv.Atoi(lines[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %w", caseIdx+1, err)
		}
		pos += 2
		if pos+2*m > len(lines) {
			return nil, fmt.Errorf("case %d: insufficient edge data", caseIdx+1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			a, err := strconv.Atoi(lines[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d: parse a: %w", caseIdx+1, i+1, err)
			}
			b, err := strconv.Atoi(lines[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d: parse b: %w", caseIdx+1, i+1, err)
			}
			edges[i] = [2]int{a, b}
			pos += 2
		}
		cases = append(cases, testCase{n: n, m: m, edges: edges})
	}
	return cases, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected := solveCase(tc.n, tc.m, tc.edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

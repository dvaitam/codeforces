package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
6 3 3
1 2
2 3
1 4
1 5
2 6
2 3 4
1 2
5 5 1
1 2
1 3
3 4
2 5
5 3 2
1 2
1 3
1 4
2 5
6 5 4
1 2
1 3
1 4
1 5
2 6
3 2 2
1 2
2 3
3 5 2
1 2
1 3
5 3 1
1 2
2 3
1 4
2 5
4 1 3
1 2
1 3
3 4
4 1 3
1 2
2 3
2 4
4 2 4
1 2
1 3
1 4
4 1 3
1 2
1 3
3 4
5 3 4
1 2
2 3
1 4
2 5
6 2 1
1 2
2 3
2 4
3 5
5 6
4 4 1
1 2
2 3
1 4
5 1 2
1 2
2 3
1 4
3 5
4 5 1
1 2
2 3
2 4
3 1 2
1 2
2 3
3 1 1
1 2
1 3
3 5 3
1 2
2 3
5 2 1
1 2
2 3
2 4
2 5
3 5 2
1 2
1 3
3 4 3
1 2
1 3
5 3 2
1 2
1 3
3 4
4 5
5 3 4
1 2
2 3
2 4
4 5
3 1 4
1 2
2 3
4 2 1
1 2
2 3
3 4
6 5 3
1 2
2 3
3 4
1 5
3 6
4 5 3
1 2
2 3
3 4
4 3 2
1 2
2 3
3 4
4 3 3
1 2
2 3
3 4
4 3 3
1 2
2 3
2 4
3 4 3
1 2
1 3
6 2 2
1 2
2 3
2 4
1 5
4 6
3 5 5
1 2
2 3
6 5 3
1 2
1 3
1 4
4 5
5 6
3 2 2
1 2
2 3
3 2 1
1 2
1 3
4 2 4
1 2
1 3
2 4
5 3 4
1 2
1 3
1 4
3 5
2 3 4
1 2
5 1 5
1 2
1 3
3 4
3 5
2 3 5
1 2
6 3 3
1 2
2 3
2 4
3 5
4 6
3 2 5
1 2
2 3
3 2 5
1 2
2 3
2 4 1
1 2
6 5 5
1 2
2 3
3 4
4 5
4 6
5 2 4
1 2
2 3
3 4
4 5
5 2 5
1 2
2 3
2 4
3 5
5 1 3
1 2
2 3
2 4
2 5
5 1 1
1 2
1 3
2 4
1 5
3 4 1
1 2
2 3
3 4 1
1 2
2 3
4 2 3
1 2
2 3
3 4
5 5 5
1 2
1 3
1 4
3 5
4 5 3
1 2
1 3
2 4
4 3 5
1 2
1 3
1 4
4 5 4
1 2
1 3
1 4
5 4 5
1 2
2 3
2 4
4 5
3 3 4
1 2
2 3
2 4 5
1 2
5 3 2
1 2
2 3
3 4
4 5
4 5 2
1 2
1 3
1 4
3 3 1
1 2
1 3
4 4 1
1 2
2 3
1 4
2 1 2
1 2
3 3 4
1 2
1 3
6 2 4
1 2
1 3
3 4
2 5
1 6
6 1 4
1 2
2 3
1 4
1 5
4 6
6 2 1
1 2
1 3
2 4
4 5
4 6
4 3 1
1 2
2 3
1 4
6 5 3
1 2
2 3
1 4
2 5
3 6
4 5 2
1 2
2 3
2 4
5 2 2
1 2
2 3
1 4
1 5
2 1 5
1 2
3 4 4
1 2
2 3
2 5 1
1 2
2 2 4
1 2
2 1 4
1 2
3 5 4
1 2
1 3
4 3 5
1 2
2 3
2 4
3 3 3
1 2
2 3
6 1 5
1 2
1 3
2 4
1 5
1 6
2 1 2
1 2
3 1 4
1 2
2 3
2 4 4
1 2
6 4 3
1 2
2 3
1 4
2 5
2 6
3 1 3
1 2
2 3
3 5 3
1 2
2 3
3 3 2
1 2
2 3
2 3 3
1 2
2 1 4
1 2
3 1 3
1 2
2 3
3 2 1
1 2
1 3
6 2 1
1 2
2 3
3 4
3 5
2 6
5 3 3
1 2
2 3
1 4
2 5
3 2 1
1 2
2 3
2 4 1
1 2
5 2 1
1 2
2 3
1 4
1 5`

func bfs(start int, adj [][]int) ([]int, int) {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	far := start
	for head := 0; head < len(q); head++ {
		v := q[head]
		if dist[v] > dist[far] {
			far = v
		}
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist, far
}

func solveCase(n int, k, c int64, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dist1, far1 := bfs(1, adj)
	distA, far2 := bfs(far1, adj)
	distB, _ := bfs(far2, adj)
	best := int64(-1 << 63)
	for v := 1; v <= n; v++ {
		ecc := distA[v]
		if distB[v] > ecc {
			ecc = distB[v]
		}
		profit := int64(ecc)*k - int64(dist1[v])*c
		if profit > best {
			best = profit
		}
	}
	return fmt.Sprintf("%d", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runExe(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	idx := 1
	for caseNum := 0; caseNum < t; caseNum++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("case %d: missing header", caseNum+1)
		}
		header := strings.Fields(lines[idx])
		idx++
		if len(header) != 3 {
			return nil, fmt.Errorf("case %d: bad header", caseNum+1)
		}
		n, err := strconv.Atoi(header[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		k, err := strconv.ParseInt(header[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad k: %w", caseNum+1, err)
		}
		c, err := strconv.ParseInt(header[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad c: %w", caseNum+1, err)
		}
		if idx+n-1 > len(lines) {
			return nil, fmt.Errorf("case %d: missing edges", caseNum+1)
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			parts := strings.Fields(lines[idx])
			idx++
			if len(parts) != 2 {
				return nil, fmt.Errorf("case %d: bad edge line", caseNum+1)
			}
			u, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge u: %w", caseNum+1, err)
			}
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad edge v: %w", caseNum+1, err)
			}
			edges[i] = [2]int{u, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", n, k, c)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(n, k, c, edges),
		})
	}
	return cases, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `2 2 2 1 -5 1 2 0
4 1 3 4 -2
5 5 5 1 -2 5 3 5 5 1 1 3 1 0 4 3 2
2 1 2 1 -5
3 3 2 1 3 3 1 5 2 1 4
5 3 3 5 -4 1 5 5 5 2 -4
4 3 2 4 2 2 1 -2 2 1 5
4 3 3 4 3 4 3 -3 4 1 -3
5 5 4 1 2 2 1 4 4 1 3 4 3 1 1 2 5
4 4 4 2 -5 1 3 4 3 2 -1 4 3 3
4 4 4 3 4 4 2 1 4 3 5 4 3 0
3 3 2 3 4 2 3 1 2 3 -3
4 1 4 1 -3
2 1 2 1 -2
3 2 3 1 1 3 2 -4
5 2 1 3 0 2 5 0
3 1 3 1 -4
2 2 2 1 5 1 2 5
5 4 1 2 -2 4 5 -1 4 5 -4 3 1 -4
4 3 3 4 1 4 3 -5 1 2 4
3 3 3 1 -3 2 1 -5 3 1 1
4 1 2 1 1
3 2 3 2 5 1 2 3
2 1 1 2 -3
5 2 2 5 -3 4 2 0
3 3 3 1 4 2 1 0 2 3 2
2 2 1 2 -1 1 2 -4
2 1 1 2 -5
2 1 2 1 -5
3 2 1 3 -1 1 2 -4
2 1 2 1 -4
3 1 3 1 5
5 3 3 4 -2 2 5 5 2 3 4
2 2 2 1 4 1 2 0
2 1 1 2 0
5 1 3 2 -5
2 2 1 2 -5 2 1 5
4 2 2 4 3 2 4 4
5 5 5 4 -1 3 2 -1 1 2 2 1 2 4 1 5 -3
2 1 2 1 5
5 2 4 1 0 4 3 0
3 2 2 3 5 2 1 -5
4 2 1 4 4 2 1 -2
5 4 2 5 -5 5 4 0 3 5 1 3 1 -1
3 1 1 3 4
3 3 3 1 -4 1 3 -4 1 2 0
3 3 1 3 -3 2 3 2 3 1 4
3 2 1 3 0 3 1 0
4 1 3 1 -5
2 2 2 1 -5 2 1 1
4 1 3 2 -3
4 4 4 1 3 2 1 2 3 1 -3 4 3 -4
5 1 4 2 3
3 2 2 3 0 2 3 -1
5 5 1 5 -2 5 1 3 4 5 3 3 2 -5 3 1 -2
4 2 2 3 -5 3 2 -3
2 1 2 1 0
3 2 3 1 2 2 3 4
4 1 1 4 -5
3 3 2 1 5 1 2 5 3 2 3
4 3 3 4 0 1 4 5 2 3 -5
4 1 1 4 1
2 1 2 1 -4
4 3 3 4 1 1 3 2 4 3 -1
3 1 1 3 3
5 2 2 3 5 1 3 -3
2 1 2 1 -4
4 2 2 1 -5 2 3 1
5 5 5 4 2 3 5 2 2 5 -4 2 4 -2 5 3 2
4 2 2 1 1 2 4 4
2 1 2 1 1
3 2 2 3 5 2 1 -4
5 5 2 3 4 4 3 -4 4 1 -4 5 1 2 4 1 5
2 2 2 1 -3 2 1 1
2 1 2 1 4
4 3 4 3 -5 3 4 -1 2 4 0
2 1 1 2 5
2 2 2 1 -2 2 1 -4
3 1 1 2 1
3 2 2 1 0 3 2 -2
2 1 1 2 3
2 2 1 2 0 1 2 -2
2 1 2 1 3
4 2 1 3 -4 4 2 5
3 3 2 3 4 3 2 -4 2 3 3
3 1 3 1 5
4 4 3 4 5 1 4 -2 2 3 -2 4 3 -5
4 3 3 4 1 3 4 -1 2 3 5
3 3 2 3 -3 3 1 -3 1 3 -1
5 5 1 2 3 3 1 3 1 4 1 1 3 0 2 3 3
2 2 1 2 5 2 1 5
4 1 1 2 -1
3 1 1 2 -5
4 3 2 4 1 1 2 5 2 3 -3
4 2 3 2 -3 1 2 4
5 1 1 5 -5
4 3 4 3 3 2 4 4 4 3 5
5 4 5 2 -5 1 4 5 5 4 -4 5 4 0
3 1 3 1 -4
3 1 2 3 -2
2 2 1 2 -5 2 1 -2
5 2 1 4 5 3 2 -3
5 4 4 5 -2 3 2 -4 4 1 -2 3 5 4
4 1 3 2 5
4 3 3 4 1 3 4 0 2 1 4
5 3 1 2 0 5 1 -4 1 2 -5
4 1 4 1 -3
3 1 1 2 -3
3 1 1 3 -2
2 2 2 1 -5 2 1 2`

type testCase struct {
	input    string
	expected string
}

type edge struct {
	to int
	w  int64
}

func solve(n, m int, edges [][3]int64) string {
	adj := make([][]edge, n+1)
	for _, e := range edges {
		a, b, d := int(e[0]), int(e[1]), e[2]
		adj[b] = append(adj[b], edge{to: a, w: d})
		adj[a] = append(adj[a], edge{to: b, w: -d})
	}
	pos := make([]int64, n+1)
	vis := make([]bool, n+1)
	queue := make([]int, 0)
	ok := true
	for i := 1; i <= n && ok; i++ {
		if !vis[i] {
			vis[i] = true
			pos[i] = 0
			queue = append(queue, i)
			for len(queue) > 0 && ok {
				u := queue[0]
				queue = queue[1:]
				for _, e := range adj[u] {
					v := e.to
					val := pos[u] + e.w
					if !vis[v] {
						vis[v] = true
						pos[v] = val
						queue = append(queue, v)
					} else if pos[v] != val {
						ok = false
						break
					}
				}
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("case %d: not enough tokens", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad m: %w", idx+1, err)
		}
		if len(parts) != 2+3*m {
			return nil, fmt.Errorf("case %d: expected %d numbers got %d", idx+1, 2+3*m, len(parts))
		}
		edges := make([][3]int64, m)
		for i := 0; i < m; i++ {
			a, err := strconv.ParseInt(parts[2+3*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad a: %w", idx+1, err)
			}
			b, err := strconv.ParseInt(parts[2+3*i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad b: %w", idx+1, err)
			}
			d, err := strconv.ParseInt(parts[2+3*i+2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad d: %w", idx+1, err)
			}
			edges[i] = [3]int64{a, b, d}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solve(n, m, edges),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: verifierH /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

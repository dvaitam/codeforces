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
8 1 2 2 3 1 4 3 5 5 6 4 7 4 8
8 1 2 2 3 2 4 2 5 5 6 2 7 3 8
3 1 2 2 3
6 1 2 2 3 1 4 1 5 3 6
5 1 2 2 3 2 4 3 5
3 1 2 2 3
3 1 2 1 3
6 1 2 1 3 1 4 1 5 1 6
3 1 2 1 3
4 1 2 2 3 2 4
6 1 2 1 3 3 4 3 5 3 6
8 1 2 1 3 3 4 3 5 4 6 5 7 5 8
6 1 2 2 3 3 4 4 5 5 6
5 1 2 1 3 3 4 3 5
4 1 2 2 3 3 4
3 1 2 2 3
3 1 2 2 3
4 1 2 1 3 1 4
6 1 2 2 3 2 4 3 5 3 6
5 1 2 2 3 3 4 4 5
4 1 2 2 3 3 4
4 1 2 2 3 2 4
4 1 2 1 3 1 4
4 1 2 2 3 2 4
3 1 2 1 3
6 1 2 1 3 2 4 4 5 2 6
6 1 2 2 3 3 4 1 5 4 6
6 1 2 1 3 2 4 2 5 2 6
6 1 2 1 3 1 4 4 5 4 6
6 1 2 1 3 1 4 1 5 1 6
6 1 2 1 3 1 4 2 5 5 6
6 1 2 1 3 3 4 1 5 1 6
6 1 2 1 3 1 4 1 5 1 6
6 1 2 1 3 1 4 2 5 2 6
6 1 2 1 3 1 4 1 5 5 6
6 1 2 1 3 1 4 4 5 4 6
6 1 2 1 3 1 4 4 5 4 6
5 1 2 1 3 1 4 1 5
5 1 2 2 3 3 4 3 5
5 1 2 2 3 2 4 2 5
5 1 2 2 3 3 4 3 5
6 1 2 1 3 2 4 1 5 5 6
6 1 2 1 3 1 4 2 5 3 6
6 1 2 2 3 3 4 4 5 5 6
5 1 2 2 3 3 4 4 5
5 1 2 1 3 3 4 3 5
5 1 2 1 3 1 4 4 5
5 1 2 1 3 1 4 3 5
5 1 2 2 3 2 4 2 5
5 1 2 2 3 2 4 2 5
5 1 2 1 3 1 4 1 5
5 1 2 2 3 2 4 2 5
6 1 2 2 3 1 4 3 5 5 6
6 1 2 1 3 1 4 1 5 1 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 6 5 5 6
6 1 2 1 3 1 4 6 5 1 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6
6 1 2 1 3 1 4 1 5 6 6`

type Edge struct {
	to  int
	idx int
}

type Node struct {
	id    int
	idx   int
	level int
}

type testCase struct {
	input    string
	expected string
}

func solve(n int, edges [][2]int) int {
	adj := make([][]Edge, n+1)
	for i, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], Edge{v, i + 1})
		adj[v] = append(adj[v], Edge{u, i + 1})
	}
	visited := make([]bool, n+1)
	queue := []Node{{1, 0, 1}}
	visited[1] = true
	ans := 1
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		for _, e := range adj[cur.id] {
			if visited[e.to] {
				continue
			}
			lvl := cur.level
			if e.idx < cur.idx {
				lvl++
			}
			if lvl > ans {
				ans = lvl
			}
			visited[e.to] = true
			queue = append(queue, Node{e.to, e.idx, lvl})
		}
	}
	return ans
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
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
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
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
			edges[i] = [2]int{u, v}
			pos += 2
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.Itoa(solve(n, edges)),
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
		fmt.Println("usage: verifierA /path/to/binary")
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

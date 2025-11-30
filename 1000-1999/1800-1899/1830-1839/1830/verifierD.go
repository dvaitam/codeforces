package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `3 1 2 2 3
6 1 2 1 3 3 4 1 5 4 6
4 1 2 1 3 3 4
5 1 2 2 3 3 4 2 5
3 1 2 2 3
7 1 2 1 3 1 4 1 5 3 6 1 7
8 1 2 2 3 3 4 4 5 4 6 4 7 6 8
8 1 2 1 3 2 4 1 5 1 6 2 7 4 8
3 1 2 2 3
8 1 2 2 3 3 4 4 5 5 6 3 7 5 8
6 1 2 1 3 2 4 1 5 3 6
6 1 2 2 3 3 4 1 5 2 6
7 1 2 2 3 1 4 1 5 4 6 6 7
5 1 2 2 3 1 4 4 5
3 1 2 2 3
5 1 2 1 3 1 4 1 5
5 1 2 2 3 3 4 2 5
2 1 2
2 1 2
2 1 2
3 1 2 2 3
6 1 2 1 3 3 4 1 5 3 6
4 1 2 1 3 2 4
5 1 2 2 3 3 4 1 5
6 1 2 2 3 3 4 2 5 3 6
5 1 2 2 3 3 4 3 5
2 1 2
6 1 2 1 3 2 4 2 5 1 6
7 1 2 2 3 2 4 3 5 5 6 6 7
4 1 2 1 3 3 4
2 1 2
4 1 2 2 3 2 4
6 1 2 1 3 2 4 2 5 3 6
8 1 2 2 3 2 4 4 5 1 6 1 7 5 8
7 1 2 2 3 3 4 2 5 3 6 2 7
4 1 2 2 3 3 4
7 1 2 1 3 3 4 3 5 3 6 6 7
8 1 2 2 3 1 4 1 5 3 6 6 7 6 8
3 1 2 2 3
3 1 2 1 3
6 1 2 2 3 3 4 2 5 3 6
4 1 2 2 3 3 4
3 1 2 2 3
6 1 2 2 3 2 4 4 5 3 6
5 1 2 1 3 2 4 2 5
3 1 2 2 3
8 1 2 1 3 1 4 4 5 5 6 3 7 5 8
4 1 2 1 3 3 4
4 1 2 1 3 1 4
2 1 2
5 1 2 1 3 2 4 1 5
3 1 2 1 3
7 1 2 2 3 1 4 1 5 3 6 2 7
4 1 2 1 3 2 4
3 1 2 1 3
6 1 2 1 3 1 4 3 5 2 6
8 1 2 2 3 3 4 4 5 1 6 3 7 2 8
4 1 2 1 3 2 4
4 1 2 1 3 1 4
2 1 2
2 1 2
5 1 2 1 3 3 4 4 5
4 1 2 2 3 1 4
4 1 2 2 3 3 4
4 1 2 2 3 1 4
4 1 2 1 3 1 4
6 1 2 2 3 1 4 2 5 1 6
4 1 2 2 3 3 4
8 1 2 2 3 1 4 3 5 4 6 6 7 3 8
5 1 2 2 3 1 4 2 5
3 1 2 1 3
8 1 2 1 3 2 4 2 5 1 6 3 7 3 8
6 1 2 1 3 2 4 1 5 5 6
8 1 2 1 3 3 4 3 5 5 6 5 7 1 8
8 1 2 2 3 1 4 2 5 4 6 1 7 5 8
2 1 2
7 1 2 1 3 1 4 3 5 1 6 1 7
8 1 2 2 3 3 4 3 5 5 6 3 7 7 8
2 1 2
8 1 2 1 3 3 4 1 5 5 6 1 7 5 8
4 1 2 1 3 1 4
3 1 2 2 3
6 1 2 2 3 2 4 4 5 3 6
6 1 2 1 3 2 4 2 5 4 6
8 1 2 2 3 3 4 4 5 2 6 6 7 4 8
3 1 2 1 3
5 1 2 2 3 3 4 2 5
3 1 2 1 3
3 1 2 1 3
8 1 2 2 3 3 4 3 5 2 6 3 7 1 8
4 1 2 2 3 1 4
3 1 2 1 3
4 1 2 1 3 2 4
7 1 2 1 3 2 4 1 5 4 6 6 7
2 1 2
8 1 2 2 3 1 4 2 5 1 6 2 7 7 8
4 1 2 1 3 2 4
3 1 2 1 3
4 1 2 1 3 2 4
3 1 2 1 3`

type testCase struct {
	input    string
	expected string
}

func solve(n int, edges [][2]int) int64 {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	color := make([]int, n+1)
	for i := range color {
		color[i] = -1
	}
	cnt := [2]int{}
	queue := []int{1}
	color[1] = 0
	cnt[0]++
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range g[v] {
			if color[to] == -1 {
				color[to] = 1 - color[v]
				cnt[color[to]]++
				queue = append(queue, to)
			}
		}
	}
	if cnt[1] > cnt[0] {
		cnt[1], cnt[0] = cnt[0], cnt[1]
	}
	n64 := int64(n)
	return n64*n64 - int64(cnt[1])
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	pos := 0
	cases := make([]testCase, 0)
	for pos < len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("bad n at token %d: %w", pos, err)
		}
		pos++
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if pos+1 >= len(fields) {
				return nil, fmt.Errorf("case %d: incomplete edge list", len(cases)+1)
			}
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad u: %w", len(cases)+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad v: %w", len(cases)+1, err)
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
			expected: strconv.FormatInt(solve(n, edges), 10),
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
		fmt.Println("usage: verifierD /path/to/binary")
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

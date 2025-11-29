package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
6 1 2 2 3 3 4 1 5 4 6
3 1 2 1 3
2 1 2
5 1 2 2 3 3 4 1 5
6 1 2 1 3 3 4 2 5 4 6
5 1 2 1 3 3 4 2 5
4 1 2 2 3 1 4
3 1 2 2 3
2 1 2
4 1 2 1 3 1 4
2 1 2
4 1 2 1 3 2 4
3 1 2 2 3
3 1 2 1 3
2 1 2
3 1 2 1 3
2 1 2
5 1 2 1 3 3 4 2 5
4 1 2 2 3 3 4
5 1 2 1 3 1 4 4 5
5 1 2 1 3 3 4 3 5
2 1 2
4 1 2 1 3 3 4
2 1 2
2 1 2
4 1 2 1 3 3 4
5 1 2 1 3 1 4 4 5
5 1 2 2 3 3 4 3 5
2 1 2
4 1 2 1 3 1 4
4 1 2 2 3 3 4
4 1 2 1 3 1 4
5 1 2 1 3 3 4 2 5
3 1 2 1 3
5 1 2 2 3 3 4 1 5
5 1 2 1 3 1 4 3 5
5 1 2 1 3 3 4 3 5
6 1 2 1 3 1 4 2 5 3 6
5 1 2 1 3 2 4 4 5
4 1 2 2 3 2 4
5 1 2 2 3 1 4 2 5
2 1 2
5 1 2 1 3 1 4 1 5
4 1 2 1 3 2 4
4 1 2 2 3 3 4
4 1 2 2 3 1 4
5 1 2 2 3 1 4 3 5
5 1 2 2 3 2 4 2 5
5 1 2 1 3 2 4 4 5
6 1 2 1 3 1 4 2 5 4 6
5 1 2 2 3 1 4 4 5
3 1 2 1 3
4 1 2 1 3 3 4
4 1 2 2 3 3 4
5 1 2 1 3 2 4 1 5
6 1 2 1 3 3 4 2 5 3 6
4 1 2 1 3 1 4
5 1 2 1 3 3 4 3 5
6 1 2 1 3 1 4 2 5 2 6
5 1 2 2 3 3 4 3 5
6 1 2 1 3 2 4 1 5 1 6
3 1 2 2 3
4 1 2 2 3 3 4
5 1 2 2 3 2 4 2 5
3 1 2 2 3
4 1 2 1 3 3 4
4 1 2 2 3 3 4
4 1 2 1 3 1 4
4 1 2 1 3 3 4
5 1 2 1 3 2 4 4 5
4 1 2 1 3 1 4
4 1 2 1 3 3 4
4 1 2 1 3 3 4
6 1 2 2 3 3 4 4 5 5 6
5 1 2 1 3 1 4 3 5
3 1 2 2 3
5 1 2 1 3 1 4 2 5
3 1 2 2 3
5 1 2 1 3 1 4 4 5
3 1 2 1 3
4 1 2 1 3 1 4
4 1 2 2 3 2 4
3 1 2 2 3
5 1 2 1 3 1 4 4 5
4 1 2 1 3 3 4
5 1 2 1 3 2 4 1 5
5 1 2 2 3 2 4 3 5
2 1 2
5 1 2 2 3 2 4 3 5
2 1 2
4 1 2 2 3 2 4
6 1 2 1 3 1 4 1 5 1 6
3 1 2 2 3
5 1 2 2 3 1 4 4 5
4 1 2 1 3 3 4
6 1 2 2 3 2 4 2 5 2 6
4 1 2 2 3 3 4
5 1 2 1 3 3 4 1 5
`

type testCase struct {
	n     int
	edges [][2]int
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

// solve replicates 1375G.go logic to compute answer.
func solve(tc testCase) int {
	n := tc.n
	adj := make([][]int, n)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	q := []int{0}
	color[0] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if color[v] == -1 {
				color[v] = color[u] ^ 1
				q = append(q, v)
			}
		}
	}
	c0, c1 := 0, 0
	for _, c := range color {
		if c == 0 {
			c0++
		} else {
			c1++
		}
	}
	if c1 < c0 {
		c0 = c1
	}
	if c0 > 0 {
		c0--
	}
	return c0
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		expect := 1 + 2*(n-1)
		if len(fields) != expect {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, expect, len(fields))
		}
		edges := make([][2]int, n-1)
		pos := 1
		for i := 0; i < n-1; i++ {
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, err
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, err
			}
			edges[i] = [2]int{u - 1, v - 1}
			pos += 2
		}
		cases = append(cases, testCase{n: n, edges: edges})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		want := solve(tc)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(want) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

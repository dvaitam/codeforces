package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesData = `
5 4 1 2 2 3 3 4 4 5
5 3 1 2 2 3 3 4 2 5
6 2 1 2 1 3 1 4 3 5 5 6
7 5 1 2 2 3 1 4 1 5 3 6 4 7
6 1 1 2 2 3 2 4 2 5 5 6
5 4 1 2 1 3 3 4 1 5
2 1 1 2
6 4 1 2 1 3 3 4 3 5 1 6
3 1 1 2 1 3
6 4 1 2 1 3 2 4 4 5 1 6
4 3 1 2 1 3 3 4
4 3 1 2 2 3 2 4
2 1 1 2
6 2 1 2 1 3 1 4 2 5 1 6
6 3 1 2 1 3 1 4 2 5 2 6
2 1 1 2
7 5 1 2 1 3 1 4 4 5 5 6 3 7
5 4 1 2 1 3 2 4 1 5
5 3 1 2 1 3 1 4 3 5
2 1 1 2
3 2 1 2 1 3
2 1 1 2
2 1 1 2
2 1 1 2
5 1 1 2 1 3 1 4 1 5
3 1 1 2 2 3
3 1 1 2 2 3
6 1 1 2 1 3 1 4 1 5 3 6
4 2 1 2 1 3 3 4
5 1 1 2 2 3 1 4 3 5
4 3 1 2 1 3 3 4
7 2 1 2 1 3 1 4 3 5 5 6 3 7
2 1 1 2
2 1 1 2
6 5 1 2 2 3 2 4 3 5 2 6
6 1 1 2 1 3 2 4 1 5 5 6
4 1 1 2 2 3 2 4
6 3 1 2 1 3 3 4 3 5 4 6
7 4 1 2 1 3 3 4 2 5 3 6 2 7
3 1 1 2 2 3
7 6 1 2 1 3 2 4 4 5 1 6 2 7
5 1 1 2 1 3 2 4 4 5
6 5 1 2 1 3 2 4 3 5 3 6
5 1 1 2 1 3 3 4 1 5
7 2 1 2 2 3 3 4 4 5 3 6 1 7
3 1 1 2 1 3
3 1 1 2 2 3
4 3 1 2 1 3 2 4
5 1 1 2 2 3 2 4 1 5
4 1 1 2 1 3 1 4
4 1 1 2 1 3 1 4
7 3 1 2 2 3 3 4 1 5 5 6 6 7
5 4 1 2 2 3 3 4 2 5
3 2 1 2 1 3
3 1 1 2 2 3
4 2 1 2 2 3 3 4
2 1 1 2
3 1 1 2 2 3
5 2 1 2 1 3 2 4 2 5
2 1 1 2
6 1 1 2 2 3 2 4 3 5 4 6
2 1 1 2
5 3 1 2 1 3 2 4 1 5
7 4 1 2 1 3 2 4 3 5 2 6 2 7
7 5 1 2 1 3 1 4 1 5 2 6 6 7
3 1 1 2 1 3
2 1 1 2
5 4 1 2 2 3 1 4 3 5
3 2 1 2 2 3
3 2 1 2 1 3
7 1 1 2 1 3 1 4 4 5 4 6 3 7
3 1 1 2 1 3
2 1 1 2
5 3 1 2 1 3 3 4 4 5
3 2 1 2 2 3
7 3 1 2 2 3 3 4 2 5 5 6 5 7
3 1 1 2 2 3
4 1 1 2 2 3 3 4
3 2 1 2 2 3
2 1 1 2
2 1 1 2
2 1 1 2
5 3 1 2 1 3 3 4 4 5
4 3 1 2 1 3 3 4
2 1 1 2
3 2 1 2 2 3
5 2 1 2 1 3 3 4 2 5
4 3 1 2 2 3 1 4
5 3 1 2 2 3 2 4 4 5
2 1 1 2
3 2 1 2 1 3
2 1 1 2
3 2 1 2 1 3
4 2 1 2 2 3 2 4
5 1 1 2 2 3 1 4 1 5
2 1 1 2
2 1 1 2
5 4 1 2 1 3 3 4 4 5
6 3 1 2 2 3 1 4 4 5 3 6
7 6 1 2 1 3 1 4 4 5 4 6 3 7
`

type testCase struct {
	n     int
	k     int
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

// solve mirrors the reference 1336A solution.
func solve(n, k int, edges [][2]int) int {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	depth := make([]int, n+1)
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}
	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			size[u] += size[v]
		}
	}
	vals := make([]int, n)
	for i := 1; i <= n; i++ {
		vals[i-1] = depth[i] - (size[i] - 1)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	ans := 0
	for i := 0; i < k; i++ {
		ans += vals[i]
	}
	return ans
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
		if len(fields) < 3 {
			return nil, fmt.Errorf("case %d invalid line", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		expectEdges := n - 1
		if len(fields) != 2+2*expectEdges {
			return nil, fmt.Errorf("case %d invalid number of values", idx+1)
		}
		edges := make([][2]int, expectEdges)
		pos := 2
		for i := 0; i < expectEdges; i++ {
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, err
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, err
			}
			edges[i] = [2]int{u, v}
			pos += 2
		}
		cases = append(cases, testCase{n: n, k: k, edges: edges})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		want := solve(tc.n, tc.k, tc.edges)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil || ans != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

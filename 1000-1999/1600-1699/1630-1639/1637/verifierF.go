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

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `
5
4 9 3 3 8
1 2
2 3
3 5
4 5
2
4 1
1 2
3
2 9 8
1 3
2 3
3
6 1 8
1 2
2 3
5
9 2 3 3 7
1 5
2 3
3 5
4 5
4
4 7 1 9
1 2
2 4
3 4
2
7 4
1 2
5
5 3 8 5 5
1 3
2 5
3 5
4 5
4
6 8 10 3
1 3
2 3
3 4
3
9 10 1
1 2
2 3
3
3 8 3
1 3
2 3
6
4 2 9 10 9 4
1 4
2 3
3 4
4 6
5 6
4
6 10 4 1
1 3
2 3
3 4
2
2 1
1 2
2
5 8
1 2
6
10 1 3 5 4 3
1 6
2 5
3 4
4 6
5 6
3
3 6 9
1 2
2 3
6
6 9 1 9 5 4
1 6
2 4
3 4
4 6
5 6
5
10 3 4 3 6
1 5
2 5
3 4
4 5
5
7 5 7 1 4
1 5
2 5
3 4
4 5
4
8 9 6 6
1 4
2 4
3 4
5
1 1 6 7 6
1 5
2 4
3 5
4 5
2
3 1
1 2
4
1 10 3 10
1 2
2 3
3 4
6
5 2 1 2 3 5
1 3
2 6
3 4
4 5
5 6
4
1 2 1 9
1 2
2 3
3 4
5
4 6 3 9 5
1 5
2 5
3 4
4 5
3
9 5 2
1 3
2 3
5
5 5 10 4 9
1 4
2 5
3 4
4 5
6
2 6 6 10 6 2
1 4
2 4
3 4
4 5
5 6
6
2 2 6 5 10 5
1 5
2 5
3 4
4 6
5 6
3
4 6 8
1 2
2 3
6
2 6 10 8 10 7
1 6
2 6
3 5
4 5
5 6
4
4 2 3 10
1 3
2 4
3 4
3
10 4 4
1 3
2 3
4
6 4 2 1
1 3
2 3
3 4
2
3 2
1 2
2
2 10
1 2
6
2 8 3 3 9 1
1 2
2 6
3 6
4 6
5 6
4
5 6 6 9
1 2
2 3
3 4
3
3 5 5
1 3
2 3
3
7 4 10
1 2
2 3
5
9 2 7 2 4
1 4
2 5
3 5
4 5
3
6 2 2
1 2
2 3
3
7 7 1
1 3
2 3
3
5 8 1
1 3
2 3
6
8 10 5 6 6 5
1 6
2 4
3 6
4 5
5 6
5
1 7 5 4 3
1 4
2 5
3 5
4 5
6
9 8 8 6 4 1
1 6
2 4
3 6
4 6
5 6
2
10 9
1 2
2
7 5
1 2
5
8 10 1 7 8
1 2
2 5
3 5
4 5
6
6 1 1 9 7 10
1 3
2 6
3 6
4 6
5 6
3
9 7 2
1 2
2 3
2
9 2
1 2
3
9 2 2
1 2
2 3
4
6 4 9 8
1 4
2 4
3 4
2
6 6
1 2
6
3 6 1 8 4 6
1 2
2 3
3 4
4 5
5 6
4
5 1 6 8
1 2
2 3
3 4
3
8 6 6
1 2
2 3
4
4 1 3 1
1 2
2 4
3 4
6
4 9 3 7 6 3
1 6
2 3
3 6
4 5
5 6
4
6 1 1 9
1 4
2 3
3 4
4
2 6 8 4
1 4
2 4
3 4
6
6 4 4 6 9 10
1 3
2 3
3 4
4 5
5 6
4
2 3 7 8
1 3
2 4
3 4
5
1 9 2 4 3
1 4
2 4
3 4
4 5
5
5 7 6 5 3
1 2
2 4
3 5
4 5
6
7 5 2 6 3 7
1 3
2 6
3 5
4 6
5 6
3
10 1 4
1 3
2 3
2
5 4
1 2
3
6 3 5
1 3
2 3
3
2 10 4
1 3
2 3
3
8 9 2
1 3
2 3
2
3 4
1 2
3
10 3 9
1 2
2 3
6
10 3 7 7 3 5
1 5
2 5
3 6
4 5
5 6
3
4 6 5
1 3
2 3
3
6 9 6
1 3
2 3
4
6 2 4 5
1 4
2 4
3 4
2
3 2
1 2
4
2 5 9 8
1 2
2 3
3 4
5
7 5 10 2 8
1 5
2 4
3 4
4 5
4
8 10 6 7
1 4
2 3
3 4
5
9 5 3 2 2
1 4
2 3
3 5
4 5
6
7 4 4 4 4 5
1 2
2 3
3 5
4 5
5 6
2
10 5
1 2
6
9 1 6 3 2 9
1 2
2 6
3 5
4 5
5 6
4
9 8 10 2
1 3
2 3
3 4
4
8 9 7 8
1 3
2 4
3 4
2
3 9
1 2
4
8 10 5 8
1 3
2 3
3 4
4
9 3 10 2
1 4
2 4
3 4
3
2 8 2
1 2
2 3
4
10 7 7 7
1 2
2 3
3 4
3
2 6 8
1 3
2 3
3
10 10 4
1 2
2 3
3
5 9 1
1 2
2 3
2
3 5
1 2
`

type testCase struct {
	n     int
	h     []int
	edges [][2]int
}

func solveCase(tc testCase) int64 {
	n := tc.n
	h := tc.h
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// Root at node with maximum height (matches accepted solution)
	maxH := -1
	root := -1
	for i := 1; i <= n; i++ {
		if h[i] > maxH {
			maxH = h[i]
			root = i
		}
	}

	var totalSum int64

	var dfs func(u, p int) int
	dfs = func(u, p int) int {
		var childVals []int
		for _, v := range adj[u] {
			if v != p {
				childVals = append(childVals, dfs(v, u))
			}
		}
		if len(childVals) == 0 {
			return h[u]
		}
		maxVal := childVals[0]
		maxIdx := 0
		for i := 1; i < len(childVals); i++ {
			if childVals[i] > maxVal {
				maxVal = childVals[i]
				maxIdx = i
			}
		}
		for i, val := range childVals {
			if i != maxIdx {
				totalSum += int64(val)
			}
		}
		if maxVal < h[u] {
			maxVal = h[u]
		}
		return maxVal
	}

	var childVals []int
	for _, v := range adj[root] {
		childVals = append(childVals, dfs(v, root))
	}

	if len(childVals) == 0 {
		totalSum += int64(h[root]) * 2
	} else if len(childVals) == 1 {
		totalSum += int64(h[root]) * 2
	} else {
		sort.Slice(childVals, func(i, j int) bool {
			return childVals[i] > childVals[j]
		})
		totalSum += int64(h[root]) * 2
		for i := 2; i < len(childVals); i++ {
			totalSum += int64(childVals[i])
		}
	}

	return totalSum
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0)
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("line %d: missing heights", i+1)
		}
		hFields := strings.Fields(lines[i])
		if len(hFields) != n {
			return nil, fmt.Errorf("line %d: expected %d heights got %d", i+1, n, len(hFields))
		}
		h := make([]int, n+1)
		for idx, f := range hFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse h[%d]: %v", i+1, idx, err)
			}
			h[idx+1] = v
		}
		i++
		edges := make([][2]int, 0, n-1)
		for e := 0; e < n-1; e++ {
			if i >= len(lines) {
				return nil, fmt.Errorf("line %d: missing edge", i+1)
			}
			parts := strings.Fields(lines[i])
			if len(parts) != 2 {
				return nil, fmt.Errorf("line %d: expected 2 numbers got %d", i+1, len(parts))
			}
			u, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse u: %v", i+1, err)
			}
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse v: %v", i+1, err)
			}
			edges = append(edges, [2]int{u, v})
			i++
		}
		cases = append(cases, testCase{n: n, h: h, edges: edges})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx := 1; idx <= tc.n; idx++ {
			if idx > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.h[idx]))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.FormatInt(expected, 10) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

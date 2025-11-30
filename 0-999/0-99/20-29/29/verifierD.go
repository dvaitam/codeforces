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

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `6 1 2 2 3 3 4 4 5 1 6 5 6
5 1 2 1 3 2 4 2 5 4 5 3
7 1 2 1 3 1 4 4 5 5 6 4 7 3 6 2 7
6 1 2 1 3 1 4 2 5 2 6 6 5 3 4
6 1 2 2 3 3 4 3 5 1 6 4 6 5
4 1 2 2 3 3 4 4
4 1 2 2 3 1 4 3 4
5 1 2 2 3 3 4 3 5 4 5
6 1 2 1 3 1 4 1 5 5 6 3 6 4 2
7 1 2 2 3 3 4 2 5 4 6 1 7 6 7 5
3 1 2 2 3 3
5 1 2 2 3 2 4 3 5 4 5
4 1 2 1 3 3 4 4 2
3 1 2 2 3 3
7 1 2 2 3 3 4 4 5 3 6 1 7 6 5 7
6 1 2 2 3 3 4 1 5 4 6 6 5
7 1 2 2 3 1 4 4 5 1 6 4 7 7 6 5 3
7 1 2 2 3 2 4 2 5 2 6 5 7 7 4 3 6
7 1 2 2 3 1 4 1 5 4 6 4 7 5 3 6 7
5 1 2 1 3 1 4 3 5 5 4 2
3 1 2 1 3 3 2
7 1 2 1 3 1 4 4 5 5 6 6 7 2 3 7
3 1 2 2 3 3
5 1 2 1 3 1 4 2 5 3 5 4
5 1 2 1 3 2 4 3 5 4 5
3 1 2 2 3 3
7 1 2 1 3 3 4 3 5 4 6 2 7 7 6 5
4 1 2 2 3 2 4 3 4
7 1 2 1 3 2 4 3 5 5 6 4 7 7 6
5 1 2 2 3 2 4 3 5 4 5
5 1 2 2 3 2 4 1 5 5 3 4
5 1 2 1 3 2 4 4 5 5 3
3 1 2 2 3 3
6 1 2 1 3 3 4 4 5 5 6 2 6
6 1 2 2 3 3 4 1 5 1 6 4 6 5
6 1 2 2 3 1 4 1 5 4 6 6 5 3
4 1 2 1 3 1 4 4 3 2
3 1 2 1 3 3 2
6 1 2 1 3 1 4 2 5 1 6 4 5 3 6
7 1 2 1 3 3 4 2 5 2 6 2 7 4 6 5 7
3 1 2 1 3 2 3
3 1 2 1 3 2 3
5 1 2 2 3 3 4 2 5 5 4
5 1 2 1 3 2 4 1 5 4 3 5
7 1 2 2 3 1 4 3 5 2 6 4 7 6 5 7
7 1 2 2 3 1 4 4 5 4 6 2 7 3 5 6 7
4 1 2 2 3 2 4 4 3
7 1 2 1 3 2 4 2 5 3 6 4 7 7 5 6
6 1 2 2 3 2 4 2 5 4 6 6 5 3
6 1 2 1 3 1 4 4 5 3 6 2 6 5
7 1 2 1 3 1 4 3 5 5 6 6 7 4 2 7
3 1 2 2 3 3
3 1 2 1 3 2 3
6 1 2 2 3 2 4 2 5 2 6 5 4 3 6
7 1 2 2 3 1 4 4 5 2 6 6 7 3 5 7
7 1 2 2 3 1 4 2 5 4 6 5 7 7 3 6
3 1 2 1 3 3 2
4 1 2 1 3 3 4 4 2
7 1 2 2 3 3 4 3 5 4 6 6 7 7 5
4 1 2 2 3 3 4 4
6 1 2 2 3 3 4 3 5 4 6 6 5
4 1 2 1 3 1 4 2 4 3
7 1 2 2 3 3 4 1 5 1 6 4 7 6 7 5
6 1 2 1 3 3 4 1 5 2 6 4 5 6
5 1 2 1 3 1 4 4 5 3 5 2
6 1 2 2 3 1 4 1 5 2 6 6 4 5 3
4 1 2 2 3 1 4 3 4
5 1 2 2 3 2 4 3 5 5 4
7 1 2 1 3 1 4 2 5 3 6 6 7 4 5 7
3 1 2 2 3 3
5 1 2 1 3 1 4 4 5 3 2 5
4 1 2 1 3 3 4 2 4
4 1 2 2 3 2 4 3 4
3 1 2 2 3 3
5 1 2 1 3 2 4 1 5 5 4 3
3 1 2 2 3 3
4 1 2 1 3 3 4 2 4
3 1 2 1 3 3 2
7 1 2 2 3 1 4 2 5 3 6 1 7 7 6 4 5
6 1 2 2 3 2 4 2 5 3 6 4 5 6
7 1 2 1 3 2 4 3 5 2 6 3 7 6 5 7 4
5 1 2 1 3 1 4 1 5 3 5 2 4
5 1 2 2 3 2 4 1 5 3 5 4
3 1 2 2 3 3
5 1 2 2 3 1 4 4 5 3 5
3 1 2 1 3 2 3
3 1 2 2 3 3
7 1 2 1 3 3 4 2 5 5 6 2 7 7 4 6
6 1 2 2 3 3 4 2 5 3 6 4 6 5
6 1 2 2 3 2 4 2 5 5 6 4 6 3
6 1 2 1 3 1 4 1 5 3 6 4 5 6 2
6 1 2 2 3 3 4 4 5 1 6 6 5
6 1 2 2 3 2 4 4 5 4 6 5 3 6
4 1 2 1 3 1 4 3 4 2
5 1 2 1 3 3 4 3 5 4 2 5
4 1 2 2 3 3 4 4
4 1 2 1 3 3 4 4 2
4 1 2 1 3 3 4 2 4
4 1 2 2 3 1 4 3 4
7 1 2 1 3 2 4 4 5 1 6 2 7 5 7 6 3`

type testCase struct {
	n      int
	edges  [][2]int
	leaves []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: empty", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if n < 1 {
			return nil, fmt.Errorf("line %d: n must be positive", idx+1)
		}
		expectedFields := 1 + 2*(n-1)
		if len(fields) < expectedFields {
			return nil, fmt.Errorf("line %d: not enough edge data", idx+1)
		}
		edges := make([][2]int, n-1)
		ptr := 1
		for i := 0; i < n-1; i++ {
			u, err := strconv.Atoi(fields[ptr])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse u%d: %w", idx+1, i+1, err)
			}
			v, err := strconv.Atoi(fields[ptr+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse v%d: %w", idx+1, i+1, err)
			}
			edges[i] = [2]int{u, v}
			ptr += 2
		}
		deg := make([]int, n+1)
		for _, e := range edges {
			deg[e[0]]++
			deg[e[1]]++
		}
		leafCount := 0
		for v := 2; v <= n; v++ {
			if deg[v] == 1 {
				leafCount++
			}
		}
		if len(fields)-expectedFields < leafCount {
			return nil, fmt.Errorf("line %d: not enough leaf values (need %d have %d)", idx+1, leafCount, len(fields)-expectedFields)
		}
		leaves := make([]int, leafCount)
		for i := 0; i < leafCount; i++ {
			v, err := strconv.Atoi(fields[expectedFields+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse leaf %d: %w", idx+1, i+1, err)
			}
			leaves[i] = v
		}
		cases = append(cases, testCase{n: n, edges: edges, leaves: leaves})
	}
	return cases, nil
}

// solveCase mirrors 29D.go.
func solveCase(tc testCase) string {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	isLeaf := make([]bool, n+1)
	leafCount := 0
	for v := 2; v <= n; v++ {
		if len(adj[v]) == 1 {
			isLeaf[v] = true
			leafCount++
		}
	}
	if leafCount != len(tc.leaves) {
		return "-1"
	}
	pos := make([]int, n+1)
	for i, v := range tc.leaves {
		pos[v] = i
	}
	parent := make([]int, n+1)
	children := make([][]int, n+1)
	stack := []int{1}
	parent[1] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			children[v] = append(children[v], u)
			stack = append(stack, u)
		}
	}
	const inf = int(1e9)
	minpos := make([]int, n+1)
	maxpos := make([]int, n+1)
	cnt := make([]int, n+1)
	for i := 1; i <= n; i++ {
		minpos[i] = inf
		maxpos[i] = -inf
	}
	order := make([]int, 0, n)
	var dfs func(int)
	dfs = func(v int) {
		for _, u := range children[v] {
			dfs(u)
		}
		order = append(order, v)
	}
	dfs(1)
	for _, v := range order {
		if isLeaf[v] {
			cnt[v] = 1
			minpos[v] = pos[v]
			maxpos[v] = pos[v]
		}
		for _, u := range children[v] {
			cnt[v] += cnt[u]
			if minpos[u] < minpos[v] {
				minpos[v] = minpos[u]
			}
			if maxpos[u] > maxpos[v] {
				maxpos[v] = maxpos[u]
			}
		}
		if cnt[v] > 0 && maxpos[v]-minpos[v]+1 != cnt[v] {
			return "-1"
		}
	}
	for v := 1; v <= n; v++ {
		childrenV := children[v]
		sort.Slice(childrenV, func(i, j int) bool {
			return minpos[childrenV[i]] < minpos[childrenV[j]]
		})
		children[v] = childrenV
	}
	res := make([]int, 0, 2*n-1)
	var build func(int)
	build = func(v int) {
		for _, u := range children[v] {
			res = append(res, u)
			build(u)
			res = append(res, v)
		}
	}
	res = append(res, 1)
	build(1)
	if len(res) != 2*n-1 {
		return "-1"
	}
	out := make([]string, len(res))
	for i, v := range res {
		out[i] = strconv.Itoa(v)
	}
	return strings.Join(out, " ")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for i, v := range tc.leaves {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expect := solveCase(tc)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

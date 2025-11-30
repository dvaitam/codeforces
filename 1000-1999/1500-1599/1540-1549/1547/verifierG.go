package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesG.txt so the verifier is self-contained.
const testcasesRaw = `
100
5 2
2 5
4 3
4 12
2 4
1 2
3 4
2 1
4 1
3 1
4 3
1 4
4 2
2 3
3 2
1 3
5 12
1 2
3 4
2 1
1 5
4 3
5 4
5 1
1 4
4 2
5 3
3 2
3 5
4 6
2 4
3 4
4 3
4 2
1 4
2 3
1 0
1 0
3 2
2 3
1 2
5 0
4 12
2 4
1 2
2 1
3 4
4 1
3 1
4 3
4 2
1 4
2 3
3 2
1 3
3 0
2 2
1 2
2 1
2 0
3 4
2 3
3 2
1 3
3 1
4 11
1 3
2 4
1 2
3 4
2 1
4 3
3 1
4 2
2 3
3 2
4 1
5 10
2 4
2 1
4 1
3 1
4 3
1 5
5 1
5 3
3 2
1 3
5 5
2 4
4 3
3 1
5 4
3 2
2 1
1 2
4 7
2 1
3 4
4 2
1 4
2 3
3 2
4 1
5 19
4 3
3 1
5 4
5 1
2 5
1 3
4 2
4 5
5 3
2 4
1 2
2 1
1 5
3 2
4 1
3 5
5 2
1 4
2 3
1 0
5 5
1 5
3 1
1 4
4 2
3 2
5 6
3 4
3 1
5 4
1 4
3 2
3 5
5 13
1 2
2 1
1 5
3 1
4 3
5 4
5 1
4 2
2 3
5 3
2 5
1 3
3 5
4 11
1 3
2 4
3 4
2 1
4 3
3 1
4 2
1 4
2 3
3 2
4 1
1 0
1 0
1 0
1 0
4 5
3 4
2 1
4 2
2 3
3 2
3 5
2 1
3 1
2 3
3 2
1 3
3 5
1 2
2 1
2 3
3 2
1 3
2 2
1 2
2 1
4 9
1 3
2 4
3 4
2 1
4 3
3 1
4 2
2 3
4 1
3 1
2 3
1 0
2 1
2 1
5 19
3 4
4 3
3 1
5 4
5 1
2 5
4 2
4 5
5 3
2 4
1 2
2 1
1 5
3 2
4 1
3 5
5 2
1 4
2 3
1 0
1 0
2 0
1 0
2 0
2 0
5 6
4 1
3 1
5 4
4 2
1 3
5 2
5 18
1 3
2 4
1 2
2 1
3 4
4 3
1 5
3 1
5 4
5 1
1 4
4 2
2 3
4 5
5 3
3 2
4 1
5 2
2 0
2 0
2 0
2 1
1 2
1 0
5 7
3 4
4 3
1 4
2 3
5 3
3 2
4 1
1 0
3 5
1 2
2 1
2 3
3 2
1 3
4 4
3 1
2 4
2 1
4 3
5 7
2 1
4 3
1 5
5 4
1 4
2 5
3 5
4 3
4 2
3 2
1 4
3 2
2 3
3 1
4 9
1 2
3 4
4 1
4 3
3 1
1 4
2 3
3 2
1 3
2 2
1 2
2 1
5 8
2 4
3 1
5 1
4 2
2 3
2 5
4 1
5 2
1 0
4 2
1 2
2 1
4 2
2 4
4 1
5 1
3 5
1 0
5 9
2 4
3 4
1 5
3 1
1 4
4 2
2 3
2 5
4 1
4 0
3 0
2 2
1 2
2 1
5 18
2 4
1 2
3 4
2 1
1 5
4 1
4 3
5 4
3 1
5 1
4 2
2 3
4 5
5 3
3 2
2 5
1 3
5 2
1 0
4 5
2 1
3 1
4 2
1 4
1 3
3 4
2 3
3 2
1 3
3 1
4 12
2 4
1 2
2 1
3 4
4 3
3 1
4 1
4 2
1 4
2 3
3 2
1 3
5 20
3 4
4 3
3 1
5 4
5 1
2 5
1 3
4 2
4 5
5 3
2 4
1 2
2 1
1 5
3 2
4 1
3 5
5 2
1 4
2 3
1 0
1 0
5 11
2 1
3 4
4 3
1 5
5 1
1 4
2 3
4 5
3 2
1 3
3 5
1 0
2 0
5 7
1 2
5 4
5 1
1 4
2 3
3 2
2 5
3 1
2 1
2 2
1 2
2 1
1 0
2 1
1 2
2 0
4 9
1 2
2 1
3 4
4 1
4 3
3 1
4 2
3 2
1 3
3 0
2 0
5 13
2 4
2 1
3 4
1 5
3 1
5 4
5 1
4 2
1 4
5 3
3 2
4 1
3 5
3 5
1 2
2 1
3 1
2 3
1 3
4 8
1 3
2 1
3 4
3 1
4 2
2 3
3 2
4 1
1 0
4 4
2 4
2 1
3 4
4 3
1 0
5 3
5 4
4 1
5 2
3 3
3 1
1 2
2 3
2 0
1 0
2 0

`

type edge struct {
	a, b int
}

type testCase struct {
	n     int
	m     int
	edges []edge
}

func solveCase(tc testCase) []int {
	n, m := tc.n, tc.m
	g := make([][]int, n+1)
	rg := make([][]int, n+1)
	for _, e := range tc.edges {
		g[e.a] = append(g[e.a], e.b)
		rg[e.b] = append(rg[e.b], e.a)
	}

	reachable := make([]bool, n+1)
	stack := []int{1}
	reachable[1] = true
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if !reachable[to] {
				reachable[to] = true
				stack = append(stack, to)
			}
		}
	}

	order := make([]int, 0, n)
	used := make([]bool, n+1)
	var dfs1 func(int)
	dfs1 = func(v int) {
		used[v] = true
		for _, to := range g[v] {
			if reachable[to] && !used[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for v := 1; v <= n; v++ {
		if reachable[v] && !used[v] {
			dfs1(v)
		}
	}

	comp := make([]int, n+1)
	compCnt := 0
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for _, to := range rg[v] {
			if reachable[to] && comp[to] == 0 {
				dfs2(to, c)
			}
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == 0 {
			compCnt++
			dfs2(v, compCnt)
		}
	}

	compSize := make([]int, compCnt+1)
	for v := 1; v <= n; v++ {
		if reachable[v] {
			compSize[comp[v]]++
		}
	}

	cyc := make([]bool, compCnt+1)
	for c := 1; c <= compCnt; c++ {
		if compSize[c] > 1 {
			cyc[c] = true
		}
	}
	for v := 1; v <= n; v++ {
		if reachable[v] {
			cv := comp[v]
			for _, to := range g[v] {
				if reachable[to] && comp[to] == cv && v == to {
					cyc[cv] = true
				}
			}
		}
	}

	adjC := make([][]int, compCnt+1)
	for v := 1; v <= n; v++ {
		if !reachable[v] {
			continue
		}
		cv := comp[v]
		for _, to := range g[v] {
			if !reachable[to] {
				continue
			}
			ct := comp[to]
			if cv != ct {
				adjC[cv] = append(adjC[cv], ct)
			}
		}
	}

	startComp := comp[1]
	reachComp := make([]bool, compCnt+1)
	queue := []int{startComp}
	reachComp[startComp] = true
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, to := range adjC[c] {
			if !reachComp[to] {
				reachComp[to] = true
				queue = append(queue, to)
			}
		}
	}

	inf := make([]bool, compCnt+1)
	q := make([]int, 0)
	for c := 1; c <= compCnt; c++ {
		if reachComp[c] && cyc[c] {
			inf[c] = true
			q = append(q, c)
		}
	}
	for head := 0; head < len(q); head++ {
		c := q[head]
		for _, to := range adjC[c] {
			if reachComp[to] && !inf[to] {
				inf[to] = true
				q = append(q, to)
			}
		}
	}

	dp := make([]int, compCnt+1)
	indeg := make([]int, compCnt+1)
	for c := 1; c <= compCnt; c++ {
		if !reachComp[c] || inf[c] {
			continue
		}
		for _, to := range adjC[c] {
			if reachComp[to] && !inf[to] {
				indeg[to]++
			}
		}
	}

	queue = queue[:0]
	for c := 1; c <= compCnt; c++ {
		if !reachComp[c] || inf[c] {
			continue
		}
		if indeg[c] == 0 {
			queue = append(queue, c)
		}
	}
	dp[startComp] = 1
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, to := range adjC[c] {
			if !reachComp[to] || inf[to] {
				continue
			}
			val := dp[to] + dp[c]
			if val > 2 {
				val = 2
			}
			if val > dp[to] {
				dp[to] = val
			}
			indeg[to]--
			if indeg[to] == 0 {
				queue = append(queue, to)
			}
		}
	}

	ans := make([]int, n)
	for v := 1; v <= n; v++ {
		if !reachable[v] {
			ans[v-1] = 0
		} else if inf[comp[v]] {
			ans[v-1] = -1
		} else {
			ans[v-1] = dp[comp[v]]
		}
	}
	_ = m // unused after building graph
	return ans
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n/m", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", i+1, err)
		}
		pos += 2
		if pos+m*2 > len(fields) {
			return nil, fmt.Errorf("case %d: not enough edge data", i+1)
		}
		tc := testCase{n: n, m: m, edges: make([]edge, m)}
		for j := 0; j < m; j++ {
			a, err := strconv.Atoi(fields[pos+j*2])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse edge %d a: %v", i+1, j, err)
			}
			b, err := strconv.Atoi(fields[pos+j*2+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse edge %d b: %v", i+1, j, err)
			}
			tc.edges[j] = edge{a: a, b: b}
		}
		pos += m * 2
		cases = append(cases, tc)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing data after parsing")
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
	return strings.TrimRight(out.String(), "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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

	var inputBuilder strings.Builder
	inputBuilder.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		inputBuilder.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, e := range tc.edges {
			inputBuilder.WriteString(fmt.Sprintf("%d %d", e.a, e.b))
			if i+1 == tc.m {
				inputBuilder.WriteByte('\n')
			} else {
				inputBuilder.WriteByte('\n')
			}
		}
	}

	gotOut, err := runCandidate(bin, inputBuilder.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimRight(gotOut, "\n"), "\n")
	if len(lines) != len(cases) {
		fmt.Printf("expected %d lines of output, got %d\n", len(cases), len(lines))
		os.Exit(1)
	}

	for i, tc := range cases {
		expected := solveCase(tc)
		got := strings.Fields(lines[i])
		if len(got) != len(expected) {
			fmt.Printf("case %d: expected %d values, got %d\n", i+1, len(expected), len(got))
			os.Exit(1)
		}
		for j, val := range got {
			v, err := strconv.Atoi(val)
			if err != nil {
				fmt.Printf("case %d: non-integer output %q\n", i+1, val)
				os.Exit(1)
			}
			if v != expected[j] {
				fmt.Printf("case %d failed at position %d\nexpected: %d\ngot: %d\n", i+1, j+1, expected[j], v)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

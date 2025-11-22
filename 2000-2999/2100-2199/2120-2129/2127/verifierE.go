package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const maxLog = 20

type graph struct {
	n    int
	g    [][]int
	tin  []int
	tout []int
	up   [][]int
	dep  []int
	tm   int
}

func newGraph(n int) *graph {
	gr := &graph{
		n:    n,
		g:    make([][]int, n),
		tin:  make([]int, n),
		tout: make([]int, n),
		up:   make([][]int, maxLog),
		dep:  make([]int, n),
	}
	for i := 0; i < maxLog; i++ {
		gr.up[i] = make([]int, n)
	}
	return gr
}

func (gr *graph) addEdge(u, v int) {
	gr.g[u] = append(gr.g[u], v)
	gr.g[v] = append(gr.g[v], u)
}

func (gr *graph) prepare(root int) {
	type node struct {
		v   int
		p   int
		out bool
	}
	stack := []node{{root, root, false}}
	gr.tm = 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.out {
			gr.tout[cur.v] = gr.tm
			continue
		}
		gr.tin[cur.v] = gr.tm
		gr.tm++
		gr.up[0][cur.v] = cur.p
		for i := 1; i < maxLog; i++ {
			gr.up[i][cur.v] = gr.up[i-1][gr.up[i-1][cur.v]]
		}
		if cur.v == cur.p {
			gr.dep[cur.v] = 0
		} else {
			gr.dep[cur.v] = gr.dep[cur.p] + 1
		}
		stack = append(stack, node{cur.v, cur.p, true})
		for i := len(gr.g[cur.v]) - 1; i >= 0; i-- {
			to := gr.g[cur.v][i]
			if to == cur.p {
				continue
			}
			stack = append(stack, node{to, cur.v, false})
		}
	}
}

func (gr *graph) isAnc(a, b int) bool {
	return gr.tin[a] <= gr.tin[b] && gr.tout[b] <= gr.tout[a]
}

func (gr *graph) lca(a, b int) int {
	if gr.isAnc(a, b) {
		return a
	}
	if gr.isAnc(b, a) {
		return b
	}
	for i := maxLog - 1; i >= 0; i-- {
		if !gr.isAnc(gr.up[i][a], b) {
			a = gr.up[i][a]
		}
	}
	return gr.up[0][a]
}

// computeCost evaluates the cost of a fully colored tree.
func computeCost(gr *graph, weights []int64, colors []int, k int) int64 {
	colorNodes := make([][]int, k+1)
	for i, col := range colors {
		colorNodes[col] = append(colorNodes[col], i)
	}

	isCutie := make([]bool, gr.n)

	for col, nodes := range colorNodes {
		if col == 0 || len(nodes) < 2 {
			continue
		}
		sort.Slice(nodes, func(i, j int) bool { return gr.tin[nodes[i]] < gr.tin[nodes[j]] })
		uniq := nodes[:0]
		prev := -1
		for _, v := range nodes {
			if v != prev {
				uniq = append(uniq, v)
				prev = v
			}
		}
		nodes = uniq
		if len(nodes) < 2 {
			continue
		}

		deg := make(map[int]int)
		stack := make([]int, 0, len(nodes))
		stack = append(stack, nodes[0])
		addEdge := func(p, ch int) {
			deg[p]++
		}
		for i := 1; i < len(nodes); i++ {
			u := nodes[i]
			l := gr.lca(u, stack[len(stack)-1])
			for len(stack) >= 2 && gr.dep[stack[len(stack)-2]] >= gr.dep[l] {
				addEdge(stack[len(stack)-2], stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if stack[len(stack)-1] != l {
				addEdge(l, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
				if len(stack) == 0 || stack[len(stack)-1] != l {
					stack = append(stack, l)
				}
			}
			stack = append(stack, u)
		}
		for len(stack) > 1 {
			addEdge(stack[len(stack)-2], stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}

		for v, d := range deg {
			if d >= 2 && colors[v] != col {
				isCutie[v] = true
			}
		}
	}

	var cost int64
	for i, c := range isCutie {
		if c {
			cost += weights[i]
		}
	}
	return cost
}

// expectedResult computes minimal cost and a valid coloring using the reference algorithm.
func expectedResult(n, k int, weights []int64, initColors []int, edges [][2]int) (int64, []int) {
	c := make([]int, n)
	copy(c, initColors)

	gr := newGraph(n)
	for _, e := range edges {
		gr.addEdge(e[0], e[1])
	}
	gr.prepare(0)

	colorNodes := make(map[int][]int)
	for i, col := range c {
		if col != 0 {
			colorNodes[col] = append(colorNodes[col], i)
		}
	}

	countMulti := make([]int, n)
	ownMulti := make([]bool, n)
	repeatColor := make([]int, n) // 0 none, -1 multiple, >0 unique

	processColor := func(col int, nodes []int) {
		if len(nodes) < 2 {
			return
		}
		sort.Slice(nodes, func(i, j int) bool { return gr.tin[nodes[i]] < gr.tin[nodes[j]] })
		uniq := make([]int, 0, len(nodes))
		prev := -1
		for _, v := range nodes {
			if v != prev {
				uniq = append(uniq, v)
				prev = v
			}
		}
		nodes = uniq
		if len(nodes) < 2 {
			return
		}

		stack := make([]int, 0, len(nodes))
		stack = append(stack, nodes[0])
		deg := make(map[int]int)

		addEdge := func(p, ch int) {
			deg[p]++
		}

		for i := 1; i < len(nodes); i++ {
			u := nodes[i]
			l := gr.lca(u, stack[len(stack)-1])
			for len(stack) >= 2 && gr.dep[stack[len(stack)-2]] >= gr.dep[l] {
				addEdge(stack[len(stack)-2], stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if stack[len(stack)-1] != l {
				addEdge(l, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
				if len(stack) == 0 || stack[len(stack)-1] != l {
					stack = append(stack, l)
				}
			}
			stack = append(stack, u)
		}
		for len(stack) > 1 {
			addEdge(stack[len(stack)-2], stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}

		for v, d := range deg {
			if d >= 2 {
				countMulti[v]++
				if c[v] == col {
					ownMulti[v] = true
				}
				if repeatColor[v] == 0 {
					repeatColor[v] = col
				} else if repeatColor[v] != col {
					repeatColor[v] = -1
				}
			}
		}
	}

	for col, nodes := range colorNodes {
		processColor(col, nodes)
	}

	var cost int64
	for i := 0; i < n; i++ {
		mult := countMulti[i]
		if c[i] != 0 {
			if ownMulti[i] {
				mult--
			}
			if mult > 0 {
				cost += weights[i]
			}
		} else {
			if mult >= 2 {
				cost += weights[i]
			}
		}
	}

	if c[0] == 0 {
		if countMulti[0] == 1 && repeatColor[0] > 0 {
			c[0] = repeatColor[0]
		} else {
			c[0] = 1
		}
	}

	type stItem struct {
		v int
		p int
	}
	stack := []stItem{{0, 0}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range gr.g[cur.v] {
			if to == cur.p {
				continue
			}
			if c[to] == 0 {
				if countMulti[to] == 1 && repeatColor[to] > 0 {
					c[to] = repeatColor[to]
				} else {
					c[to] = c[cur.v]
				}
			}
			stack = append(stack, stItem{to, cur.v})
		}
	}

	return cost, c
}

type test struct {
	raw string
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(2127))
	var tests []test

	// Small deterministic case.
	{
		var sb strings.Builder
		sb.WriteString("2\n")
		sb.WriteString("4 4\n5 5 5 5\n1 0 2 3\n1 2\n1 3\n1 4\n")
		sb.WriteString("3 3\n3 4 5\n0 0 0\n1 2\n2 3\n")
		tests = append(tests, test{raw: sb.String()})
	}

	for len(tests) < 80 {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			n := rng.Intn(12) + 3
			k := rng.Intn(n-1) + 2
			fmt.Fprintf(&sb, "%d %d\n", n, k)
			for i := 0; i < n; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", rng.Intn(10)+1)
			}
			sb.WriteByte('\n')
			for i := 0; i < n; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				if rng.Intn(4) == 0 {
					sb.WriteByte('0')
				} else {
					fmt.Fprintf(&sb, "%d", rng.Intn(k)+1)
				}
			}
			sb.WriteByte('\n')
			for i := 2; i <= n; i++ {
				p := rng.Intn(i-1) + 1
				fmt.Fprintf(&sb, "%d %d\n", p, i)
			}
		}
		tests = append(tests, test{raw: sb.String()})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseAndJudge(rawInput, output string) error {
	in := bufio.NewReader(strings.NewReader(rawInput))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t")
	}

	type testCase struct {
		n, k int
		w    []int64
		c    []int
		e    [][2]int
	}
	tc := make([]testCase, t)
	for idx := 0; idx < t; idx++ {
		var n, k int
		fmt.Fscan(in, &n, &k)
		w := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &w[i])
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		edges := make([][2]int, 0, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			edges = append(edges, [2]int{u, v})
		}
		tc[idx] = testCase{n: n, k: k, w: w, c: c, e: edges}
	}

	expCost := make([]int64, t)
	for i := 0; i < t; i++ {
		cost, _ := expectedResult(tc[i].n, tc[i].k, tc[i].w, tc[i].c, tc[i].e)
		expCost[i] = cost
	}

	out := bufio.NewReader(strings.NewReader(output))
	for idx, cur := range tc {
		var declared int64
		if _, err := fmt.Fscan(out, &declared); err != nil {
			return fmt.Errorf("test %d: missing declared cost", idx+1)
		}
		resColors := make([]int, cur.n)
		for i := 0; i < cur.n; i++ {
			if _, err := fmt.Fscan(out, &resColors[i]); err != nil {
				return fmt.Errorf("test %d: missing color %d", idx+1, i+1)
			}
		}

		for i := 0; i < cur.n; i++ {
			if cur.c[i] != 0 && resColors[i] != cur.c[i] {
				return fmt.Errorf("test %d: color %d changed fixed value (%d->%d)", idx+1, i+1, cur.c[i], resColors[i])
			}
			if resColors[i] < 1 || resColors[i] > cur.k {
				return fmt.Errorf("test %d: color %d out of range", idx+1, i+1)
			}
		}

		gr := newGraph(cur.n)
		for _, e := range cur.e {
			gr.addEdge(e[0], e[1])
		}
		gr.prepare(0)

		actualCost := computeCost(gr, cur.w, resColors, cur.k)
		if actualCost != declared {
			return fmt.Errorf("test %d: declared cost %d but computed %d", idx+1, declared, actualCost)
		}
		if actualCost != expCost[idx] {
			return fmt.Errorf("test %d: cost %d not minimal (expected %d)", idx+1, actualCost, expCost[idx])
		}
	}

	var extra string
	if _, err := fmt.Fscan(out, &extra); err == nil {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.raw)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := parseAndJudge(tc.raw, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, tc.raw, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

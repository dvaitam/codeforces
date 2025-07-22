package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type edge struct {
	to int
	w  int
}

type tree struct {
	n      int
	adj    [][]edge
	parent []int
	weight []int
	depth  []int
	order  []int
}

func newTree(n int) *tree {
	t := &tree{n: n}
	t.adj = make([][]edge, n+1)
	t.parent = make([]int, n+1)
	t.weight = make([]int, n+1)
	t.depth = make([]int, n+1)
	t.order = make([]int, n+1)
	return t
}

func (t *tree) add(u, v, w int) {
	t.adj[u] = append(t.adj[u], edge{v, w})
	t.adj[v] = append(t.adj[v], edge{u, w})
}

func (t *tree) dfs() {
	visited := make([]bool, t.n+1)
	stack := []int{1}
	order := 0
	t.parent[1] = 0
	t.weight[1] = 0
	t.depth[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[v] {
			continue
		}
		visited[v] = true
		t.order[v] = order
		order++
		for _, e := range t.adj[v] {
			if e.to == t.parent[v] {
				continue
			}
			t.parent[e.to] = v
			t.weight[e.to] = e.w
			t.depth[e.to] = t.depth[v] + 1
			stack = append(stack, e.to)
		}
	}
}

func (t *tree) pathEdges(u, v int, set map[[2]int]bool) {
	for t.depth[u] > t.depth[v] {
		k := edgeKey(u, t.parent[u])
		set[k] = true
		u = t.parent[u]
	}
	for t.depth[v] > t.depth[u] {
		k := edgeKey(v, t.parent[v])
		set[k] = true
		v = t.parent[v]
	}
	for u != v {
		k1 := edgeKey(u, t.parent[u])
		set[k1] = true
		u = t.parent[u]
		k2 := edgeKey(v, t.parent[v])
		set[k2] = true
		v = t.parent[v]
	}
}

func edgeKey(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}

func minimalSubtree(t *tree, active []int) int {
	if len(active) <= 1 {
		return 0
	}
	set := make(map[[2]int]bool)
	for i := 0; i < len(active); i++ {
		for j := i + 1; j < len(active); j++ {
			t.pathEdges(active[i], active[j], set)
		}
	}
	sum := 0
	for e := range set {
		// find weight
		u, v := e[0], e[1]
		if t.parent[u] == v {
			sum += t.weight[u]
		} else {
			sum += t.weight[v]
		}
	}
	return sum
}

func expected(n int, edges [][3]int, ops []string) []int {
	tr := newTree(n)
	for _, e := range edges {
		tr.add(e[0], e[1], e[2])
	}
	tr.dfs()
	active := make(map[int]bool)
	answers := []int{}
	for _, op := range ops {
		if op[0] == '+' {
			var x int
			fmt.Sscanf(op[1:], "%d", &x)
			active[x] = true
		} else if op[0] == '-' {
			var x int
			fmt.Sscanf(op[1:], "%d", &x)
			delete(active, x)
		} else if op[0] == '?' {
			list := make([]int, 0, len(active))
			for v := range active {
				list = append(list, v)
			}
			sort.Slice(list, func(i, j int) bool { return tr.order[list[i]] < tr.order[list[j]] })
			ans := minimalSubtree(tr, list)
			answers = append(answers, ans)
		}
	}
	return answers
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rand.Seed(4)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(5) + 2
		edges := make([][3]int, 0, n-1)
		// generate random tree using union method
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			w := rand.Intn(10) + 1
			edges = append(edges, [3]int{p, i, w})
		}
		q := rand.Intn(6) + 3
		ops := make([]string, q)
		active := make(map[int]bool)
		ensuredQuestion := false
		for i := 0; i < q; i++ {
			typ := rand.Intn(3)
			switch typ {
			case 0:
				x := rand.Intn(n) + 1
				ops[i] = fmt.Sprintf("+ %d", x)
				active[x] = true
			case 1:
				if len(active) == 0 {
					x := rand.Intn(n) + 1
					ops[i] = fmt.Sprintf("+ %d", x)
					active[x] = true
				} else {
					var list []int
					for v := range active {
						list = append(list, v)
					}
					x := list[rand.Intn(len(list))]
					ops[i] = fmt.Sprintf("- %d", x)
					delete(active, x)
				}
			default:
				ops[i] = "?"
				ensuredQuestion = true
			}
		}
		if !ensuredQuestion {
			ops[q-1] = "?"
		}
		// Build input
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&in, "%d %d %d\n", e[0], e[1], e[2])
		}
		fmt.Fprintf(&in, "%d\n", q)
		for _, op := range ops {
			fmt.Fprintf(&in, "%s\n", op)
		}
		input := in.String()
		exp := expected(n, edges, ops)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tcase+1, err)
			return
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != len(exp) {
			fmt.Printf("test %d failed: expected %d lines got %d\ninput:\n%s", tcase+1, len(exp), len(lines), input)
			return
		}
		for i, line := range lines {
			var got int
			fmt.Sscan(line, &got)
			if got != exp[i] {
				fmt.Printf("test %d failed on answer %d: expected %d got %d\ninput:\n%s", tcase+1, i+1, exp[i], got, input)
				return
			}
		}
	}
	fmt.Println("All tests passed.")
}

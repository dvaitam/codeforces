package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
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

const MAXLOG = 19

type Node struct {
	a, b, g int
}

func (nd *Node) bad() bool { return nd.a == -1 }
func (nd *Node) beBad()    { nd.a, nd.b, nd.g = -1, -1, -1 }

func solveC(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	where := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
		where[a[i]] = i
	}
	v := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x int
		fmt.Fscan(r, &x)
		v[x-1] = append(v[x-1], i+1)
	}
	tin := make([]int, n)
	tout := make([]int, n)
	gl := make([]int, n)
	par := make([][]int, n)
	for i := range par {
		par[i] = make([]int, MAXLOG)
	}
	curT := 0
	var dfs func(int, int)
	dfs = func(cur, p int) {
		curT++
		tin[cur] = curT
		par[cur][0] = p
		for i := 1; i < MAXLOG; i++ {
			par[cur][i] = par[par[cur][i-1]][i-1]
		}
		for _, to := range v[cur] {
			gl[to] = gl[cur] + 1
			dfs(to, cur)
		}
		tout[cur] = curT
	}
	dfs(0, 0)
	isAncestor := func(x, y int) bool { return tin[x] <= tin[y] && tin[y] <= tout[x] }
	onTheWay := func(x, y, z int) bool { return isAncestor(x, y) && isAncestor(y, z) }
	var lca func(int, int) int
	lca = func(x, y int) int {
		if isAncestor(x, y) {
			return x
		}
		if isAncestor(y, x) {
			return y
		}
		for i := MAXLOG - 1; i >= 0; i-- {
			p := par[x][i]
			if !isAncestor(p, y) {
				x = p
			}
		}
		return par[x][0]
	}
	add := func(nd *Node, x int) {
		if nd.bad() {
			return
		}
		if nd.a == nd.b {
			nd.b = x
			nd.g = lca(nd.a, nd.b)
			if nd.a == nd.g {
				nd.a, nd.b = nd.b, nd.a
			}
		} else {
			if isAncestor(nd.a, x) {
				nd.a = x
				return
			}
			if nd.b == nd.g {
				if onTheWay(nd.b, x, nd.a) {
					return
				}
				if onTheWay(x, nd.b, nd.a) {
					nd.b = x
					nd.g = x
					return
				}
				ng := lca(nd.a, x)
				if isAncestor(ng, nd.b) {
					nd.b = x
					nd.g = ng
					return
				}
				nd.beBad()
				return
			}
			if isAncestor(nd.b, x) {
				nd.b = x
				return
			}
			if onTheWay(nd.g, x, nd.a) || onTheWay(nd.g, x, nd.b) {
				return
			}
			nd.beBad()
		}
	}
	merge := func(s, t Node) Node {
		if s.bad() || t.bad() {
			return Node{-1, -1, -1}
		}
		res := s
		add(&res, t.a)
		add(&res, t.b)
		return res
	}
	tree := make([]Node, 4*n)
	var build func(int, int, int)
	build = func(cur, l, r int) {
		if l == r {
			tree[cur] = Node{where[l], where[l], where[l]}
		} else {
			m := (l + r) >> 1
			lc := cur << 1
			build(lc, l, m)
			build(lc|1, m+1, r)
			tree[cur] = merge(tree[lc], tree[lc|1])
		}
	}
	build(1, 0, n-1)
	var update func(int, int, int, int)
	update = func(cur, l, r, pos int) {
		if l == r {
			tree[cur] = Node{where[l], where[l], where[l]}
		} else {
			m := (l + r) >> 1
			lc := cur << 1
			if pos <= m {
				update(lc, l, m, pos)
			} else {
				update(lc|1, m+1, r, pos)
			}
			tree[cur] = merge(tree[lc], tree[lc|1])
		}
	}
	ans := Node{where[0], where[0], where[0]}
	mex := 0
	var getAns func(int, int, int)
	getAns = func(cur, l, r int) {
		res := merge(ans, tree[cur])
		if !res.bad() {
			mex = r
			ans = res
		} else if l < r {
			m := (l + r) >> 1
			lc := cur << 1
			getAns(lc, l, m)
			if mex == m {
				getAns(lc|1, m+1, r)
			}
		}
	}
	var q int
	fmt.Fscan(r, &q)
	var sb strings.Builder
	for q > 0 {
		q--
		var t int
		fmt.Fscan(r, &t)
		if t == 1 {
			var x, y int
			fmt.Fscan(r, &x, &y)
			x--
			y--
			a[x], a[y] = a[y], a[x]
			where[a[x]] = x
			where[a[y]] = y
			update(1, 0, n-1, a[x])
			update(1, 0, n-1, a[y])
		} else {
			ans = Node{where[0], where[0], where[0]}
			mex = 0
			getAns(1, 0, n-1)
			fmt.Fprintf(&sb, "%d\n", mex+1)
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		perm := rand.Perm(n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", perm[i])
		}
		sb.WriteByte('\n')
		for i := 1; i < n; i++ {
			p := rng.Intn(i) + 1
			fmt.Fprintf(&sb, "%d ", p)
		}
		sb.WriteByte('\n')
		q := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				x := rng.Intn(n) + 1
				y := rng.Intn(n) + 1
				fmt.Fprintf(&sb, "1 %d %d\n", x, y)
			} else {
				fmt.Fprintf(&sb, "2\n")
			}
		}
		input := sb.String()
		expected := solveC(input)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

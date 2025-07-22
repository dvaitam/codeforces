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

type Edge struct{ to, id int }

var (
	n, m      int
	g         [][]Edge
	disc, low []int
	timer     int
	isBridge  []bool
	comp      []int
	compCount int
	tree      [][]int
	up        [][]int
	depth     []int
	LOG       int
)

func dfs(u, pid int) {
	timer++
	disc[u] = timer
	low[u] = timer
	for _, e := range g[u] {
		v, id := e.to, e.id
		if id == pid {
			continue
		}
		if disc[v] == 0 {
			dfs(v, id)
			if low[v] < low[u] {
				low[u] = low[v]
			}
			if low[v] > disc[u] {
				isBridge[id] = true
			}
		} else if disc[v] < low[u] {
			low[u] = disc[v]
		}
	}
}

func dfs2(u int) {
	comp[u] = compCount
	for _, e := range g[u] {
		if isBridge[e.id] {
			continue
		}
		v := e.to
		if comp[v] == 0 {
			dfs2(v)
		}
	}
}

func dfs3(u, p int) {
	up[u][0] = p
	for i := 1; i < LOG; i++ {
		up[u][i] = up[up[u][i-1]][i-1]
	}
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs3(v, u)
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := 0; i < LOG; i++ {
		if diff&(1<<i) != 0 {
			u = up[u][i]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[u][i] != up[v][i] {
			u = up[u][i]
			v = up[v][i]
		}
	}
	return up[u][0]
}

func solveB(input string) string {
	timer = 0
	reader := bufio.NewReader(strings.NewReader(input))
	var k int
	fmt.Fscan(reader, &n, &m)
	g = make([][]Edge, n+1)
	for i := 1; i <= m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		g[a] = append(g[a], Edge{b, i})
		g[b] = append(g[b], Edge{a, i})
	}
	disc = make([]int, n+1)
	low = make([]int, n+1)
	isBridge = make([]bool, m+1)
	for i := 1; i <= n; i++ {
		if disc[i] == 0 {
			dfs(i, -1)
		}
	}
	comp = make([]int, n+1)
	compCount = 0
	for i := 1; i <= n; i++ {
		if comp[i] == 0 {
			compCount++
			dfs2(i)
		}
	}
	tree = make([][]int, compCount+1)
	for u := 1; u <= n; u++ {
		for _, e := range g[u] {
			v, id := e.to, e.id
			if isBridge[id] && comp[u] < comp[v] {
				cu, cv := comp[u], comp[v]
				tree[cu] = append(tree[cu], cv)
				tree[cv] = append(tree[cv], cu)
			}
		}
	}
	for LOG = 1; (1 << LOG) <= compCount; LOG++ {
	}
	up = make([][]int, compCount+1)
	depth = make([]int, compCount+1)
	for i := range up {
		up[i] = make([]int, LOG)
	}
	for i := 1; i <= compCount; i++ {
		if up[i][0] == 0 {
			dfs3(i, i)
		}
	}
	fmt.Fscan(reader, &k)
	var out strings.Builder
	for i := 0; i < k; i++ {
		var s, l int
		fmt.Fscan(reader, &s, &l)
		u, v := comp[s], comp[l]
		w := lca(u, v)
		dist := depth[u] + depth[v] - 2*depth[w]
		out.WriteString(fmt.Sprintf("%d\n", dist))
	}
	return strings.TrimSpace(out.String())
}

type testB struct{ input, expect string }

func genTests() []testB {
	rand.Seed(42)
	var tests []testB
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 2
		maxM := n * (n - 1) / 2
		m := rand.Intn(maxM) + 1
		edges := make(map[[2]int]struct{})
		for len(edges) < m {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			edges[[2]int{a, b}] = struct{}{}
		}
		var pairs [][2]int
		for e := range edges {
			pairs = append(pairs, [2]int{e[0], e[1]})
		}
		k := rand.Intn(5) + 1
		var qs [][2]int
		for j := 0; j < k; j++ {
			s := rand.Intn(n) + 1
			l := rand.Intn(n) + 1
			qs = append(qs, [2]int{s, l})
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(pairs)))
		for _, e := range pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for _, q := range qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		input := sb.String()
		expect := solveB(input)
		tests = append(tests, testB{input, expect})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

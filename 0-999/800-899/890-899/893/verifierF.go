package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type SegTree struct {
	n    int
	tree []int64
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	t := &SegTree{n: size, tree: make([]int64, size<<1)}
	const INF = int64(1 << 60)
	for i := range t.tree {
		t.tree[i] = INF
	}
	return t
}

func (t *SegTree) Update(pos int, val int64) {
	idx := pos + t.n
	t.tree[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		if t.tree[idx<<1] < t.tree[idx<<1|1] {
			t.tree[idx] = t.tree[idx<<1]
		} else {
			t.tree[idx] = t.tree[idx<<1|1]
		}
	}
}

func (t *SegTree) Query(l, r int) int64 {
	const INF = int64(1 << 60)
	res := INF
	l += t.n
	r += t.n + 1
	for l < r {
		if l&1 == 1 {
			if t.tree[l] < res {
				res = t.tree[l]
			}
			l++
		}
		if r&1 == 1 {
			r--
			if t.tree[r] < res {
				res = t.tree[r]
			}
		}
		l >>= 1
		r >>= 1
	}
	return res
}

var (
	g       [][]int
	tin     []int
	tout    []int
	depth   []int
	order   int
	byDepth [][]int
)

func dfs(u, p, d int) {
	depth[u] = d
	if d >= len(byDepth) {
		tmp := make([][]int, d+1)
		copy(tmp, byDepth)
		byDepth = tmp
	}
	byDepth[d] = append(byDepth[d], u)
	order++
	tin[u] = order
	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u, d+1)
	}
	tout[u] = order
}

func solveCase(n, rroot int, a []int64, edges [][2]int, queries [][2]int) []int64 {
	g = make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	tin = make([]int, n+1)
	tout = make([]int, n+1)
	depth = make([]int, n+1)
	byDepth = make([][]int, 1)
	order = 0
	dfs(rroot, 0, 0)
	maxDepth := len(byDepth) - 1

	seg := NewSegTree(n + 2)
	curDepth := -1
	last := 0
	res := make([]int64, len(queries))
	for i, pq := range queries {
		p, q := pq[0], pq[1]
		x := (p+last)%n + 1
		k := (q + last) % n
		t := depth[x] + k
		if t > maxDepth {
			t = maxDepth
		}
		for curDepth < t {
			curDepth++
			for _, v := range byDepth[curDepth] {
				seg.Update(tin[v], a[v])
			}
		}
		ans := seg.Query(tin[x], tout[x])
		res[i] = ans
		last = int(ans)
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	r := rng.Intn(n) + 1
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(rng.Intn(100) + 1)
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	m := rng.Intn(10) + 1
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		queries[i][0] = rng.Intn(n)
		queries[i][1] = rng.Intn(n)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, r))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for _, pq := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", pq[0], pq[1]))
	}
	res := solveCase(n, r, a, edges, queries)
	answers := make([]string, m)
	for i := range res {
		answers[i] = fmt.Sprint(res[i])
	}
	return sb.String(), strings.Join(answers, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

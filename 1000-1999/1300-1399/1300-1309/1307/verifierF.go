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

type DSU struct {
	p []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) {
	rx, ry := d.Find(x), d.Find(y)
	if rx != ry {
		d.p[ry] = rx
	}
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func solve(input string) string {
	in := strings.NewReader(input)
	var n, k, r int
	fmt.Fscan(in, &n, &k, &r)
	adj := make([][]int, n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		edges[i] = [2]int{u, v}
	}
	rest := make([]int, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(in, &rest[i])
		rest[i]--
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	for _, u := range rest {
		dist[u] = 0
		q = append(q, u)
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	dsu := NewDSU(n)
	for _, e := range edges {
		u, v := e[0], e[1]
		if dist[u]+dist[v]+1 <= k {
			dsu.Union(u, v)
		}
	}
	compID := make([]int, n)
	idMap := make(map[int]int, n)
	cid := 0
	for i := 0; i < n; i++ {
		r := dsu.Find(i)
		if _, ok := idMap[r]; !ok {
			idMap[r] = cid
			cid++
		}
		compID[i] = idMap[r]
	}
	C := cid
	cadj := make([][]int, C)
	for _, e := range edges {
		u, v := e[0], e[1]
		cu, cv := compID[u], compID[v]
		if cu != cv {
			cadj[cu] = append(cadj[cu], cv)
			cadj[cv] = append(cadj[cv], cu)
		}
	}
	LOG := 1
	for (1 << LOG) <= C {
		LOG++
	}
	up := make([][]int, LOG)
	depth := make([]int, C)
	for i := range up {
		up[i] = make([]int, C)
		for j := range up[i] {
			up[i][j] = -1
		}
	}
	var dfs func(u, p int)
	dfs = func(u, p int) {
		up[0][u] = p
		for _, v := range cadj[u] {
			if v == p {
				continue
			}
			depth[v] = depth[u] + 1
			dfs(v, u)
		}
	}
	for i := 0; i < C; i++ {
		if up[0][i] == -1 {
			dfs(i, -1)
		}
	}
	for j := 1; j < LOG; j++ {
		for i := 0; i < C; i++ {
			if up[j-1][i] < 0 {
				up[j][i] = -1
			} else {
				up[j][i] = up[j-1][up[j-1][i]]
			}
		}
	}
	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for j := 0; j < LOG; j++ {
			if diff>>j&1 == 1 {
				a = up[j][a]
			}
		}
		if a == b {
			return a
		}
		for j := LOG - 1; j >= 0; j-- {
			if up[j][a] != up[j][b] {
				a = up[j][a]
				b = up[j][b]
			}
		}
		return up[0][a]
	}
	var vq int
	fmt.Fscan(in, &vq)
	var out strings.Builder
	for i := 0; i < vq; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		ca, cb := compID[a], compID[b]
		w := lca(ca, cb)
		distc := depth[ca] + depth[cb] - 2*depth[w]
		if distc <= k {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	return strings.TrimSpace(out.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	k := rng.Intn(3) + 1
	r := rng.Intn(n) + 1
	edges := make([][2]int, n-1)
	adj := make(map[[2]int]bool)
	for i := 1; i < n; i++ {
		p := rng.Intn(i) + 1
		edges[i-1] = [2]int{i + 1, p}
		adj[[2]int{i + 1, p}] = true
		adj[[2]int{p, i + 1}] = true
	}
	rest := rng.Perm(n)[:r]
	for i := range rest {
		rest[i]++
	}
	vq := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, r)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for i, v := range rest {
		if i+1 == r {
			fmt.Fprintf(&sb, "%d\n", v)
		} else {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	fmt.Fprintf(&sb, "%d\n", vq)
	for i := 0; i < vq; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for b == a {
			b = rng.Intn(n) + 1
		}
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

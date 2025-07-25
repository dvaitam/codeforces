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

func solveCase(input string) string {
	in := strings.NewReader(input)
	var n int
	fmt.Fscan(in, &n)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	const LOG = 18
	parent := make([][]int, LOG)
	for i := range parent {
		parent[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	order := make([]int, 0, n)
	q := []int{1}
	parent[0][1] = 0
	depth[1] = 0
	for idx := 0; idx < len(q); idx++ {
		u := q[idx]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[0][u] {
				continue
			}
			parent[0][v] = u
			depth[v] = depth[u] + 1
			q = append(q, v)
		}
	}
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		sz[i] = 1
	}
	for i := len(order) - 1; i > 0; i-- {
		u := order[i]
		p := parent[0][u]
		sz[p] += sz[u]
	}
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			parent[k][v] = parent[k-1][parent[k-1][v]]
		}
	}
	lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		dd := depth[u] - depth[v]
		for k := 0; k < LOG; k++ {
			if dd>>k&1 == 1 {
				u = parent[k][u]
			}
		}
		if u == v {
			return u
		}
		for k := LOG - 1; k >= 0; k-- {
			if parent[k][u] != parent[k][v] {
				u = parent[k][u]
				v = parent[k][v]
			}
		}
		return parent[0][u]
	}
	jump := func(u, d int) int {
		for k := 0; k < LOG && u != 0; k++ {
			if d>>k&1 == 1 {
				u = parent[k][u]
			}
		}
		return u
	}
	var m int
	fmt.Fscan(in, &m)
	var out strings.Builder
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x == y {
			fmt.Fprintf(&out, "%d\n", n)
			continue
		}
		w := lca(x, y)
		d := depth[x] + depth[y] - 2*depth[w]
		if d&1 == 1 {
			fmt.Fprintf(&out, "0\n")
			continue
		}
		k := d / 2
		var c int
		dx := depth[x] - depth[w]
		if k <= dx {
			c = jump(x, k)
		} else {
			rem := k - dx
			c = jump(y, depth[y]-depth[w]-rem)
		}
		var nx int
		if lca(c, x) == c {
			nx = jump(x, depth[x]-depth[c]-1)
		} else {
			nx = parent[0][c]
		}
		var ny int
		if lca(c, y) == c {
			ny = jump(y, depth[y]-depth[c]-1)
		} else {
			ny = parent[0][c]
		}
		var sx, sy int
		if parent[0][nx] == c {
			sx = sz[nx]
		} else {
			sx = n - sz[c]
		}
		if parent[0][ny] == c {
			sy = sz[ny]
		} else {
			sy = n - sz[c]
		}
		ans := n - sx - sy
		fmt.Fprintf(&out, "%d\n", ans)
	}
	return strings.TrimRight(out.String(), "\n")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	m := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		inp := generateCase(rng)
		exp := solveCase(inp)
		if err := runCase(bin, inp, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

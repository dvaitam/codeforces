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

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	restStops := make([]int, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(in, &restStops[i])
	}

	// BFS from rest stops
	D := make([]int, n+1)
	for i := 1; i <= n; i++ {
		D[i] = -1
	}
	queue := make([]int, 0, n)
	for _, rs := range restStops {
		D[rs] = 0
		queue = append(queue, rs)
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if D[v] == -1 {
				D[v] = D[u] + 1
				queue = append(queue, v)
			}
		}
	}

	// DSU
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}
	union := func(i, j int) {
		ri := find(i)
		rj := find(j)
		if ri != rj {
			parent[ri] = rj
		}
	}

	// LCA setup via DFS
	depth := make([]int, n+1)
	up := make([][20]int, n+1)

	var dfs func(u, p, d int)
	dfs = func(u, p, d int) {
		depth[u] = d
		up[u][0] = p
		for i := 1; i < 20; i++ {
			up[u][i] = up[up[u][i-1]][i-1]
		}
		for _, v := range adj[u] {
			if v != p {
				dfs(v, u, d+1)
				if D[u]+D[v]+1 <= k {
					union(u, v)
				}
			}
		}
	}
	dfs(1, 0, 1)

	getLCA := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		for i := 19; i >= 0; i-- {
			if depth[u]-(1<<i) >= depth[v] {
				u = up[u][i]
			}
		}
		if u == v {
			return u
		}
		for i := 19; i >= 0; i-- {
			if up[u][i] != up[v][i] {
				u = up[u][i]
				v = up[v][i]
			}
		}
		return up[u][0]
	}

	getDist := func(u, v int) int {
		return depth[u] + depth[v] - 2*depth[getLCA(u, v)]
	}

	var vq int
	fmt.Fscan(in, &vq)
	var out strings.Builder
	for qi := 0; qi < vq; qi++ {
		var a, b int
		fmt.Fscan(in, &a, &b)

		if getDist(a, b) <= k {
			out.WriteString("YES\n")
			continue
		}
		if D[a] > k || D[b] > k {
			out.WriteString("NO\n")
			continue
		}

		lca := getLCA(a, b)
		dist_a_lca := depth[a] - depth[lca]
		var X int
		if dist_a_lca+D[lca] <= k {
			if getDist(a, b)+D[b] <= k {
				X = b
			} else {
				u := b
				for j := 19; j >= 0; j-- {
					nxt := up[u][j]
					if nxt != 0 && depth[nxt] > depth[lca] {
						d_a_nxt := dist_a_lca + (depth[nxt] - depth[lca])
						if d_a_nxt+D[nxt] > k {
							u = nxt
						}
					}
				}
				X = up[u][0]
			}
		} else {
			u := a
			for j := 19; j >= 0; j-- {
				nxt := up[u][j]
				if nxt != 0 && depth[nxt] >= depth[lca] {
					d_a_nxt := depth[a] - depth[nxt]
					if d_a_nxt+D[nxt] <= k {
						u = nxt
					}
				}
			}
			X = u
		}

		dist_b_lca := depth[b] - depth[lca]
		var Y int
		if dist_b_lca+D[lca] <= k {
			if getDist(b, a)+D[a] <= k {
				Y = a
			} else {
				u := a
				for j := 19; j >= 0; j-- {
					nxt := up[u][j]
					if nxt != 0 && depth[nxt] > depth[lca] {
						d_b_nxt := dist_b_lca + (depth[nxt] - depth[lca])
						if d_b_nxt+D[nxt] > k {
							u = nxt
						}
					}
				}
				Y = up[u][0]
			}
		} else {
			u := b
			for j := 19; j >= 0; j-- {
				nxt := up[u][j]
				if nxt != 0 && depth[nxt] >= depth[lca] {
					d_b_nxt := depth[b] - depth[nxt]
					if d_b_nxt+D[nxt] <= k {
						u = nxt
					}
				}
			}
			Y = u
		}

		if getDist(a, X) >= getDist(a, Y) {
			out.WriteString("YES\n")
		} else if find(X) == find(Y) {
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
	for i := 1; i < n; i++ {
		p := rng.Intn(i) + 1
		edges[i-1] = [2]int{i + 1, p}
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

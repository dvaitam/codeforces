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

const LOG = 20

type Edge struct {
	to int
	w  int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dfs(v, p, w int, adj [][]Edge, up, mx [][]int, depth []int) {
	up[v][0] = p
	mx[v][0] = w
	for i := 1; i < LOG; i++ {
		up[v][i] = up[up[v][i-1]][i-1]
		if up[v][i] == 0 {
			mx[v][i] = mx[v][i-1]
		} else {
			mx[v][i] = max(mx[v][i-1], mx[up[v][i-1]][i-1])
		}
	}
	for _, e := range adj[v] {
		if e.to == p {
			continue
		}
		depth[e.to] = depth[v] + 1
		dfs(e.to, v, e.w, adj, up, mx, depth)
	}
}

func maxEdge(u, v int, up, mx [][]int, depth []int) int {
	if u == v {
		return 0
	}
	res := 0
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff>>i&1 == 1 {
			if mx[u][i] > res {
				res = mx[u][i]
			}
			u = up[u][i]
		}
	}
	if u == v {
		return res
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[u][i] != up[v][i] {
			if mx[u][i] > res {
				res = mx[u][i]
			}
			if mx[v][i] > res {
				res = mx[v][i]
			}
			u = up[u][i]
			v = up[v][i]
		}
	}
	if mx[u][0] > res {
		res = mx[u][0]
	}
	if mx[v][0] > res {
		res = mx[v][0]
	}
	return res
}

func solve(n int, edges [][3]int, queries [][]int) []int {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
	}
	up := make([][]int, n+1)
	mx := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		up[i] = make([]int, LOG)
		mx[i] = make([]int, LOG)
	}
	depth := make([]int, n+1)
	dfs(1, 1, 0, adj, up, mx, depth)
	open := make([]bool, n+1)
	var ans []int
	for _, q := range queries {
		if q[0] == 1 {
			l, r := q[1], q[2]
			for i := l; i <= r; i++ {
				open[i] = true
			}
		} else if q[0] == 2 {
			l, r := q[1], q[2]
			for i := l; i <= r; i++ {
				open[i] = false
			}
		} else {
			x := q[1]
			best := -1
			for i := 1; i <= n; i++ {
				if open[i] && i != x {
					w := maxEdge(x, i, up, mx, depth)
					if w > best {
						best = w
					}
				}
			}
			ans = append(ans, best)
		}
	}
	return ans
}

func buildInput(n int, edges [][3]int, queries [][]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(queries)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	for _, q := range queries {
		if q[0] == 1 || q[0] == 2 {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q[0], q[1], q[2]))
		} else {
			sb.WriteString(fmt.Sprintf("3 %d\n", q[1]))
		}
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		edges := make([][3]int, n-1)
		for v := 2; v <= n; v++ {
			u := rng.Intn(v-1) + 1
			w := rng.Intn(10) + 1
			edges[v-2] = [3]int{u, v, w}
		}
		q := rng.Intn(5) + 1
		queries := make([][]int, q)
		for j := 0; j < q; j++ {
			t := rng.Intn(3) + 1
			if t == 1 || t == 2 {
				l := rng.Intn(n) + 1
				r := rng.Intn(n) + 1
				if l > r {
					l, r = r, l
				}
				queries[j] = []int{t, l, r}
			} else {
				x := rng.Intn(n) + 1
				queries[j] = []int{3, x}
			}
		}
		input := buildInput(n, edges, queries)
		res := solve(n, edges, queries)
		var exp strings.Builder
		for idx, v := range res {
			if idx > 0 {
				exp.WriteByte('\n')
			}
			exp.WriteString(fmt.Sprintf("%d", v))
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp.String() {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp.String(), out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---------- embedded solver (from cf_t23_1389_G.go) ----------

func solveEmbedded(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)

	var n, m, k int
	if _, err := fmt.Fscanf(reader, "%d %d %d\n", &n, &m, &k); err != nil {
		return ""
	}

	special := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscanf(reader, "%d", &special[i])
	}
	fmt.Fscanf(reader, "\n")

	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscanf(reader, "%d", &c[i])
	}
	fmt.Fscanf(reader, "\n")

	w := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscanf(reader, "%d", &w[i])
	}
	fmt.Fscanf(reader, "\n")

	type Edge struct {
		to, id int
	}
	adj := make([][]Edge, n+1)
	for i := 1; i <= m; i++ {
		var u, v int
		fmt.Fscanf(reader, "%d %d\n", &u, &v)
		adj[u] = append(adj[u], Edge{v, i})
		adj[v] = append(adj[v], Edge{u, i})
	}

	// Find biconnected components using iterative Tarjan
	dfn := make([]int, n+1)
	low := make([]int, n+1)
	st := make([]int, 0)
	inSt := make([]bool, n+1)
	timer := 0
	scc := make([]int, n+1)
	sccCnt := 0

	var tarjan func(u, p int)
	tarjan = func(u, p int) {
		timer++
		dfn[u] = timer
		low[u] = timer
		st = append(st, u)
		inSt[u] = true
		for _, e := range adj[u] {
			if e.to == p {
				continue
			}
			if dfn[e.to] != 0 {
				if inSt[e.to] && dfn[e.to] < low[u] {
					low[u] = dfn[e.to]
				}
			} else {
				tarjan(e.to, u)
				if low[e.to] < low[u] {
					low[u] = low[e.to]
				}
			}
		}
		if low[u] == dfn[u] {
			sccCnt++
			for {
				v := st[len(st)-1]
				st = st[:len(st)-1]
				inSt[v] = false
				scc[v] = sccCnt
				if v == u {
					break
				}
			}
		}
	}
	tarjan(1, 0)

	treeC := make([]int64, sccCnt+1)
	for i := 1; i <= n; i++ {
		treeC[scc[i]] += c[i]
	}

	type TreeEdge struct {
		to int
		w  int64
	}
	tree := make([][]TreeEdge, sccCnt+1)
	for u := 1; u <= n; u++ {
		for _, e := range adj[u] {
			if scc[u] != scc[e.to] {
				tree[scc[u]] = append(tree[scc[u]], TreeEdge{scc[e.to], w[e.id]})
			}
		}
	}

	dpDown := make([]int64, sccCnt+1)
	var dfsDown func(u, p int)
	dfsDown = func(u, p int) {
		dpDown[u] = treeC[u]
		for _, e := range tree[u] {
			if e.to == p {
				continue
			}
			dfsDown(e.to, u)
			val := dpDown[e.to] - e.w
			if val > 0 {
				dpDown[u] += val
			}
		}
	}
	dfsDown(1, 0)

	dpUp := make([]int64, sccCnt+1)
	var dfsUp func(u, p int, upVal int64)
	dfsUp = func(u, p int, upVal int64) {
		dpUp[u] = upVal
		for _, e := range tree[u] {
			if e.to == p {
				continue
			}
			curVal := dpDown[u]
			if dpDown[e.to]-e.w > 0 {
				curVal -= (dpDown[e.to] - e.w)
			}
			if upVal > 0 {
				curVal += upVal
			}
			nextUp := curVal - e.w
			if nextUp < 0 {
				nextUp = 0
			}
			dfsUp(e.to, u, nextUp)
		}
	}
	dfsUp(1, 0, 0)

	ans := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		u := scc[i]
		ans[i] = dpDown[u] + dpUp[u]
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
	writer.Flush()
	return strings.TrimSpace(out.String())
}

// ---------- verifier infrastructure ----------

type VEdge struct {
	u, v int
	w    int64
}

func bfs(adj [][]int, start int) []bool {
	n := len(adj)
	vis := make([]bool, n)
	q := []int{start}
	vis[start] = true
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if !vis[to] {
				vis[to] = true
				q = append(q, to)
			}
		}
	}
	return vis
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func genGraph(rng *rand.Rand, n, m int) []VEdge {
	edges := make([]VEdge, 0, m)
	// ensure tree edges for connectivity
	for i := 1; i < n; i++ {
		u := i - 1
		v := i
		w := int64(rng.Intn(5))
		edges = append(edges, VEdge{u, v, w})
	}
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		dup := false
		for _, e := range edges {
			if (e.u == u && e.v == v) || (e.u == v && e.v == u) {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		w := int64(rng.Intn(5))
		edges = append(edges, VEdge{u, v, w})
	}
	return edges
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(minInt(maxEdges, 4-n+1)) + n - 1
	k := rng.Intn(n) + 1
	specials := make([]int, k)
	used := make([]bool, n)
	for i := 0; i < k; i++ {
		for {
			v := rng.Intn(n)
			if !used[v] {
				specials[i] = v
				used[v] = true
				break
			}
		}
	}
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		c[i] = int64(rng.Intn(6))
	}
	edges := genGraph(rng, n, m)
	w := make([]int64, m)
	for i := 0; i < m; i++ {
		w[i] = int64(rng.Intn(6))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", specials[i]+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", w[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", edges[i].u+1, edges[i].v+1))
	}

	// Use embedded solver as reference
	exp := solveEmbedded(sb.String())
	return sb.String(), exp
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

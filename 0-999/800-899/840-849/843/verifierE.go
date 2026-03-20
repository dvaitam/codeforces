package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e15

type Edge struct {
	to, rev   int
	cap, flow int
}

func addEdge(adj [][]Edge, u, v, cap int) {
	adj[u] = append(adj[u], Edge{v, len(adj[v]), cap, 0})
	adj[v] = append(adj[v], Edge{u, len(adj[u]) - 1, 0, 0})
}

func bfs(adj [][]Edge, s, t int, level []int) bool {
	for i := range level {
		level[i] = -1
	}
	level[s] = 0
	q := []int{s}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, e := range adj[u] {
			if e.cap > e.flow && level[e.to] == -1 {
				level[e.to] = level[u] + 1
				q = append(q, e.to)
			}
		}
	}
	return level[t] != -1
}

func dfs(adj [][]Edge, u, t int, pushed int, level, ptr []int) int {
	if pushed == 0 || u == t {
		return pushed
	}
	for ptr[u] < len(adj[u]) {
		e := &adj[u][ptr[u]]
		tr := e.to
		if level[u]+1 != level[tr] || e.cap <= e.flow {
			ptr[u]++
			continue
		}
		push := pushed
		if e.cap-e.flow < push {
			push = e.cap - e.flow
		}
		flow := dfs(adj, tr, t, push, level, ptr)
		if flow == 0 {
			ptr[u]++
			continue
		}
		e.flow += flow
		adj[tr][e.rev].flow -= flow
		return flow
	}
	return 0
}

func dinic(adj [][]Edge, s, t int) int {
	flow := 0
	level := make([]int, len(adj))
	ptr := make([]int, len(adj))
	for bfs(adj, s, t, level) {
		for i := range ptr {
			ptr[i] = 0
		}
		for {
			pushed := dfs(adj, s, t, INF, level, ptr)
			if pushed == 0 {
				break
			}
			flow += pushed
		}
	}
	return flow
}

type InputEdge struct {
	u, v, g int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, s, t int
	if _, err := fmt.Fscan(reader, &n, &m, &s, &t); err != nil {
		return
	}

	edges := make([]InputEdge, m)
	adj1 := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].g)
		if edges[i].g == 1 {
			addEdge(adj1, edges[i].u, edges[i].v, 1)
			addEdge(adj1, edges[i].v, edges[i].u, INF)
		} else {
			addEdge(adj1, edges[i].u, edges[i].v, INF)
		}
	}

	k := dinic(adj1, s, t)

	visited := make([]bool, n+1)
	q := []int{s}
	visited[s] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, e := range adj1[u] {
			if e.cap > e.flow && !visited[e.to] {
				visited[e.to] = true
				q = append(q, e.to)
			}
		}
	}

	adj2 := make([][]Edge, n+2)
	D := make([]int, n+1)
	edgeIdx := make([]int, m)
	for i := 0; i < m; i++ {
		if edges[i].g == 1 {
			u, v := edges[i].u, edges[i].v
			D[v]++
			D[u]--
			edgeIdx[i] = len(adj2[u])
			addEdge(adj2, u, v, INF)
		}
	}

	SS := 0
	TT := n + 1
	for i := 1; i <= n; i++ {
		if D[i] > 0 {
			addEdge(adj2, SS, i, D[i])
		} else if D[i] < 0 {
			addEdge(adj2, i, TT, -D[i])
		}
	}
	addEdge(adj2, t, s, INF)

	dinic(adj2, SS, TT)

	fmt.Fprintln(writer, k)
	for i := 0; i < m; i++ {
		u, v := edges[i].u, edges[i].v
		f_i := 0
		if edges[i].g == 1 {
			f_i = 1 + adj2[u][edgeIdx[i]].flow
		}

		c_i := f_i + 1
		if visited[u] && !visited[v] {
			c_i = f_i
		}
		fmt.Fprintf(writer, "%d %d\n", f_i, c_i)
	}
}
`

func buildEmbeddedRef(dir string) (string, error) {
	src := filepath.Join(dir, "ref_embedded_843E.go")
	if err := os.WriteFile(src, []byte(refSource), 0644); err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "ref_843E_bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref: %v\n%s", err, out)
	}
	os.Remove(src)
	return bin, nil
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(5) + 1
	s := 1
	t := n
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, s, t))
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		for u == t {
			u = rng.Intn(n) + 1
		}
		v := rng.Intn(n) + 1
		for v == s || v == u {
			v = rng.Intn(n) + 1
		}
		g := rng.Intn(2)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, g))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tmpDir, err := os.MkdirTemp("", "v843E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	ref, err := buildEmbeddedRef(tmpDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		exp, err := runBin(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBin(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

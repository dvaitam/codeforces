package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const LOG = 20

var (
	n     int
	adj   [][]int
	up    [LOG + 1][]int
	depth []int
	value []int64
)

func dfs(root int) {
	stack := []int{root}
	parent := make([]int, n+1)
	parent[root] = 0
	depth[root] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	up[0] = parent
	for k := 1; k <= LOG; k++ {
		up[k] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := LOG; k >= 0; k-- {
		if diff&(1<<uint(k)) != 0 {
			a = up[k][a]
		}
	}
	if a == b {
		return a
	}
	for k := LOG; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func getPath(u, v int) []int {
	w := lca(u, v)
	var path []int
	x := u
	for x != w {
		path = append(path, x)
		x = up[0][x]
	}
	path = append(path, w)
	var tmp []int
	x = v
	for x != w {
		tmp = append(tmp, x)
		x = up[0][x]
	}
	for i := len(tmp) - 1; i >= 0; i-- {
		path = append(path, tmp[i])
	}
	return path
}

func update(u, v, k, d int) {
	path := getPath(u, v)
	visited := make([]bool, n+1)
	q := make([]int, 0, len(path))
	dist := make([]int, 0, len(path))
	for _, x := range path {
		if !visited[x] {
			visited[x] = true
			q = append(q, x)
			dist = append(dist, 0)
		}
	}
	head := 0
	for head < len(q) {
		node := q[head]
		dd := dist[head]
		head++
		value[node] += int64(k)
		if dd == d {
			continue
		}
		for _, to := range adj[node] {
			if !visited[to] {
				visited[to] = true
				q = append(q, to)
				dist = append(dist, dd+1)
			}
		}
	}
}

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func oracle(input string) string {
	r := strings.NewReader(input)
	fmt.Fscan(r, &n)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(r, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	depth = make([]int, n+1)
	for i := 0; i <= LOG; i++ {
		up[i] = make([]int, n+1)
	}
	value = make([]int64, n+1)
	dfs(1)
	var m int
	fmt.Fscan(r, &m)
	var out strings.Builder
	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(r, &t)
		if t == 1 {
			var v int
			fmt.Fscan(r, &v)
			fmt.Fprintf(&out, "%d\n", value[v])
		} else {
			var u, v, k, d int
			fmt.Fscan(r, &u, &v, &k, &d)
			update(u, v, k, d)
		}
	}
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(7) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	q := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", q)
	hasPrint := false
	for i := 0; i < q; i++ {
		if rng.Intn(3) == 0 {
			v := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "1 %d\n", v)
			hasPrint = true
		} else {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			k := rng.Intn(10) + 1
			d := rng.Intn(3)
			fmt.Fprintf(&sb, "2 %d %d %d %d\n", u, v, k, d)
		}
	}
	if !hasPrint {
		fmt.Fprintf(&sb, "1 %d\n", 1)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(47))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expect := oracle(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

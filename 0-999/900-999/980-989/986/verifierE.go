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

const MOD int64 = 1000000007

type query struct{ u, v, x int }

type edge struct{ u, v int }

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func findPath(u, v int, adj [][]int) []int {
	n := len(adj)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = -1
	}
	q := []int{u}
	parent[u] = u
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur == v {
			break
		}
		for _, to := range adj[cur] {
			if parent[to] == -1 {
				parent[to] = cur
				q = append(q, to)
			}
		}
	}
	path := []int{v}
	for path[len(path)-1] != u {
		path = append(path, parent[path[len(path)-1]])
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func solveCase(n int, edges []edge, a []int, qs []query) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var sb strings.Builder
	for _, q := range qs {
		p := findPath(q.u-1, q.v-1, adj)
		res := int64(1)
		for _, node := range p {
			res = res * int64(gcd(q.x, a[node])) % MOD
		}
		fmt.Fprintf(&sb, "%d\n", res)
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(20) + 1
	}
	qcnt := rng.Intn(5) + 1
	qs := make([]query, qcnt)
	for i := 0; i < qcnt; i++ {
		qs[i] = query{rng.Intn(n) + 1, rng.Intn(n) + 1, rng.Intn(20) + 1}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", qcnt)
	for _, qu := range qs {
		fmt.Fprintf(&sb, "%d %d %d\n", qu.u, qu.v, qu.x)
	}
	input := sb.String()
	expected := solveCase(n, edges, a, qs)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

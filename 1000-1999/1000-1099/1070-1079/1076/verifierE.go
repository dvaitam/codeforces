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

func bfs(n int, adj [][]int, start, maxD int, add int64, res []int64) {
	type node struct{ v, d int }
	q := []node{{start, 0}}
	vis := make([]bool, n+1)
	vis[start] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		res[cur.v] += add
		if cur.d == maxD {
			continue
		}
		for _, to := range adj[cur.v] {
			if !vis[to] {
				vis[to] = true
				q = append(q, node{to, cur.d + 1})
			}
		}
	}
}

func solve(n int, edges [][2]int, queries [][3]int64) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	res := make([]int64, n+1)
	for _, q := range queries {
		v := int(q[0])
		d := int(q[1])
		x := q[2]
		bfs(n, adj, v, d, x, res)
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", res[i]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	m := rng.Intn(5) + 1
	queries := make([][3]int64, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		v := rng.Intn(n) + 1
		d := rng.Intn(3)
		x := rng.Intn(10) + 1
		queries[i] = [3]int64{int64(v), int64(d), int64(x)}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", v, d, x))
	}
	expect := solve(n, edges, queries)
	return sb.String(), expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expect, got)
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
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

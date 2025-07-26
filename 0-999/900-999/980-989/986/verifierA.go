package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type edge struct{ u, v int }

func solveCase(n, m, k, s int, types []int, edges []edge) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e.u-1, e.v-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dist := make([][]int, k)
	for t := 0; t < k; t++ {
		dist[t] = make([]int, n)
		for i := 0; i < n; i++ {
			dist[t][i] = -1
		}
		q := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if types[i] == t+1 {
				dist[t][i] = 0
				q = append(q, i)
			}
		}
		for head := 0; head < len(q); head++ {
			v := q[head]
			nd := dist[t][v] + 1
			for _, to := range adj[v] {
				if dist[t][to] == -1 {
					dist[t][to] = nd
					q = append(q, to)
				}
			}
		}
	}
	ans := make([]int, n)
	tmp := make([]int, k)
	for i := 0; i < n; i++ {
		for t := 0; t < k; t++ {
			tmp[t] = dist[t][i]
		}
		sort.Ints(tmp)
		sum := 0
		for j := 0; j < s; j++ {
			sum += tmp[j]
		}
		ans[i] = sum
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", ans[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(4) + 1
	s := rng.Intn(k) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	types := make([]int, n)
	for i := 0; i < n; i++ {
		types[i] = rng.Intn(k) + 1
	}
	edgeMap := make(map[[2]int]bool)
	edges := make([]edge, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if edgeMap[key] {
			continue
		}
		edgeMap[key] = true
		edges = append(edges, edge{u, v})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, k, s)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", types[i])
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	input := sb.String()
	expected := solveCase(n, m, k, s, types, edges)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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

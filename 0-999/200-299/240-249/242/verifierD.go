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

type edge struct{ u, v int }

func solve(n int, edges []edge, a []int) string {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	pt := make([]int, n)
	queue := make([]int, 0, n)
	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if pt[i] == a[i] {
			pt[i]++
			queue = append(queue, i)
			res = append(res, i)
		}
	}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, u := range adj[v] {
			pt[u]++
			if pt[u] == a[u] {
				pt[u]++
				queue = append(queue, u)
				res = append(res, u)
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(res))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	// generate simple undirected graph
	possible := make([]edge, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			possible = append(possible, edge{i, j})
		}
	}
	m := rng.Intn(len(possible) + 1)
	rng.Shuffle(len(possible), func(i, j int) { possible[i], possible[j] = possible[j], possible[i] })
	edges := possible[:m]
	a := make([]int, n)
	deg := make([]int, n)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(deg[i] + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u+1, e.v+1)
	}
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	out := solve(n, edges, a)
	return sb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

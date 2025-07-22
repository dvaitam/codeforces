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

func bfs(adj [][]int, start int) (int, []int) {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	var far int
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		far = v
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return far, dist
}

func solveF(ops []int) []int {
	n := 4
	adj := make([][]int, 2*len(ops)+5)
	for i := 2; i <= 4; i++ {
		adj[1] = append(adj[1], i)
		adj[i] = append(adj[i], 1)
	}
	res := make([]int, len(ops))
	for i, v := range ops {
		n++
		u := n
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		n++
		w := n
		adj[w] = append(adj[w], v)
		adj[v] = append(adj[v], w)
		// compute diameter by two BFS
		far, _ := bfs(adj[:n+1], 1)
		far2, dist := bfs(adj[:n+1], far)
		res[i] = dist[far2]
	}
	return res
}

func runCase(exe string, ops []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(ops))
	for _, v := range ops {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(ops) {
		return fmt.Errorf("expected %d numbers got %d", len(ops), len(fields))
	}
	expect := solveF(ops)
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("parse error: %v", err)
		}
		if v != expect[i] {
			return fmt.Errorf("at %d expected %d got %d", i+1, expect[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][]int{{1}, {2, 3}, {1, 1, 1}}
	for i := 0; i < 100; i++ {
		q := rng.Intn(5) + 1
		ops := make([]int, q)
		for j := 0; j < q; j++ {
			ops[j] = rng.Intn(j+4) + 1
		}
		cases = append(cases, ops)
	}
	for idx, ops := range cases {
		if err := runCase(exe, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

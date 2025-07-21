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

type edge struct{ a, b int }

type graphCase struct {
	n      int
	edges  []edge
	expect int
}

func compute(n int, edges []edge) int {
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e.a-1, e.b-1
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	maxProfit := 0
	dist := make([]int, n)
	parent := make([]int, n)
	skip := make([]bool, n)
	visited := make([]bool, n)
	queue := make([]int, 0, n)
	for u := 0; u < n; u++ {
		for i := 0; i < n; i++ {
			dist[i] = -1
			parent[i] = -1
		}
		queue = queue[:0]
		dist[u] = 0
		queue = append(queue, u)
		for qi := 0; qi < len(queue); qi++ {
			v := queue[qi]
			for _, w := range adj[v] {
				if dist[w] == -1 {
					dist[w] = dist[v] + 1
					parent[w] = v
					queue = append(queue, w)
				}
			}
		}
		for v := u + 1; v < n; v++ {
			pathLen := dist[v]
			for i := 0; i < n; i++ {
				skip[i] = false
				visited[i] = false
			}
			x := v
			for x != -1 {
				skip[x] = true
				x = parent[x]
			}
			best2 := 0
			for i := 0; i < n; i++ {
				if skip[i] || visited[i] {
					continue
				}
				q1 := []int{i}
				visited[i] = true
				dmap := map[int]int{i: 0}
				far := i
				for qi := 0; qi < len(q1); qi++ {
					v1 := q1[qi]
					for _, w1 := range adj[v1] {
						if skip[w1] || visited[w1] {
							continue
						}
						visited[w1] = true
						dmap[w1] = dmap[v1] + 1
						q1 = append(q1, w1)
						if dmap[w1] > dmap[far] {
							far = w1
						}
					}
				}
				q2 := []int{far}
				seen2 := map[int]bool{far: true}
				d2 := map[int]int{far: 0}
				d2max := 0
				for qi := 0; qi < len(q2); qi++ {
					v2 := q2[qi]
					for _, w2 := range adj[v2] {
						if skip[w2] || seen2[w2] {
							continue
						}
						seen2[w2] = true
						d2[w2] = d2[v2] + 1
						q2 = append(q2, w2)
						if d2[w2] > d2max {
							d2max = d2[w2]
						}
					}
				}
				if d2max > best2 {
					best2 = d2max
				}
			}
			profit := pathLen * best2
			if profit > maxProfit {
				maxProfit = profit
			}
		}
	}
	return maxProfit
}

func generateCase(rng *rand.Rand) graphCase {
	n := rng.Intn(5) + 2
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{i, p})
	}
	expect := compute(n, edges)
	return graphCase{n, edges, expect}
}

func runCase(bin string, c graphCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for _, e := range c.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != c.expect {
		return fmt.Errorf("expected %d got %d", c.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type edge struct{ to, id int }

var (
	adj   [][]edge
	disc  []int
	low   []int
	sz    []int
	timer int
	n     int
	ans   int64
)

func dfs(v, pe int) {
	timer++
	disc[v] = timer
	low[v] = timer
	sz[v] = 1
	for _, e := range adj[v] {
		if e.id == pe {
			continue
		}
		to := e.to
		if disc[to] == 0 {
			dfs(to, e.id)
			sz[v] += sz[to]
			if low[to] > disc[v] {
				s := sz[to]
				val := int64(s*(s-1)/2 + (n-s)*(n-s-1)/2)
				if val < ans {
					ans = val
				}
			}
			if low[to] < low[v] {
				low[v] = low[to]
			}
		} else {
			if disc[to] < low[v] {
				low[v] = disc[to]
			}
		}
	}
}

func solveCase(nv int, edges [][2]int) int64 {
	n = nv
	adj = make([][]edge, n)
	id := 0
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], edge{v, id})
		adj[v] = append(adj[v], edge{u, id})
		id++
	}
	disc = make([]int, n)
	low = make([]int, n)
	sz = make([]int, n)
	timer = 0
	ans = int64(n * (n - 1) / 2)
	dfs(0, -1)
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	m := n - 1 + rng.Intn(n) // up to around 2n-1
	edges := make([][2]int, 0, m)
	// build tree first
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{p, i})
	}
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		// avoid duplicates simple
		edges = append(edges, [2]int{u, v})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	expected := fmt.Sprintf("%d\n", solveCase(n, edges))
	return sb.String(), expected
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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

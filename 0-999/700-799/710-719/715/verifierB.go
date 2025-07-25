package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Edge struct{ u, v, w int }

type Case struct {
	n, m, L int
	s, t    int
	edges   []Edge
}

func shortestPath(n int, edges []Edge, s, t int) int {
	const INF = int(1e9)
	g := make([][]Edge, n)
	for _, e := range edges {
		g[e.u] = append(g[e.u], Edge{e.v, 0, e.w})
		g[e.v] = append(g[e.v], Edge{e.u, 0, e.w})
	}
	dist := make([]int, n)
	vis := make([]bool, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[s] = 0
	for {
		u := -1
		for i := 0; i < n; i++ {
			if !vis[i] && (u == -1 || dist[i] < dist[u]) {
				u = i
			}
		}
		if u == -1 || dist[u] == INF {
			break
		}
		vis[u] = true
		for _, e := range g[u] {
			if dist[e.u] > dist[u]+e.w {
				dist[e.u] = dist[u] + e.w
			}
		}
	}
	return dist[t]
}

func existsAssignment(c Case) bool {
	var zeros []int
	for i, e := range c.edges {
		if e.w == 0 {
			zeros = append(zeros, i)
		}
	}
	maxW := c.L
	if maxW < 1 {
		maxW = 1
	}
	var dfs func(int) bool
	dfs = func(idx int) bool {
		if idx == len(zeros) {
			return shortestPath(c.n, c.edges, c.s, c.t) == c.L
		}
		for w := 1; w <= maxW; w++ {
			c.edges[zeros[idx]].w = w
			if dfs(idx + 1) {
				return true
			}
		}
		return false
	}
	return dfs(0)
}

func generateCase(rng *rand.Rand) Case {
	n := rng.Intn(3) + 2 //2..4
	edges := make([]Edge, 0)
	for i := 0; i < n-1; i++ {
		w := rng.Intn(3) + 1
		edges = append(edges, Edge{i, i + 1, w})
	}
	// maybe add extra edge
	if rng.Intn(2) == 0 {
		u := rng.Intn(n)
		v := rng.Intn(n)
		for v == u {
			v = rng.Intn(n)
		}
		w := rng.Intn(3) + 1
		edges = append(edges, Edge{u, v, w})
	}
	// at most 2 zero edges
	zeroCount := 0
	for i := range edges {
		if rng.Intn(2) == 0 && zeroCount < 2 {
			edges[i].w = 0
			zeroCount++
		}
	}
	m := len(edges)
	L := rng.Intn(8) + 2
	s := 0
	t := n - 1
	return Case{n, m, L, s, t, edges}
}

func runCase(bin string, c Case, expectPossible bool) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", c.n, c.m, c.L, c.s, c.t)
	for _, e := range c.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	tokens := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	tokens.Split(bufio.ScanWords)
	if !tokens.Scan() {
		return fmt.Errorf("no output")
	}
	ans := strings.ToUpper(tokens.Text())
	if ans == "NO" {
		if expectPossible {
			return fmt.Errorf("expected YES, got NO")
		}
		if tokens.Scan() {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("first token should be YES or NO")
	}
	edgesOut := make([]Edge, c.m)
	for i := 0; i < c.m; i++ {
		if !tokens.Scan() {
			return fmt.Errorf("missing edge data")
		}
		u, _ := strconv.Atoi(tokens.Text())
		if !tokens.Scan() {
			return fmt.Errorf("missing edge data")
		}
		v, _ := strconv.Atoi(tokens.Text())
		if !tokens.Scan() {
			return fmt.Errorf("missing edge data")
		}
		w, err := strconv.Atoi(tokens.Text())
		if err != nil {
			return fmt.Errorf("bad weight")
		}
		if u != c.edges[i].u || v != c.edges[i].v {
			return fmt.Errorf("edge %d endpoints mismatch", i+1)
		}
		if c.edges[i].w > 0 && w != c.edges[i].w {
			return fmt.Errorf("edge %d weight changed", i+1)
		}
		if c.edges[i].w == 0 && w <= 0 {
			return fmt.Errorf("edge %d weight not positive", i+1)
		}
		edgesOut[i] = Edge{u, v, w}
	}
	if tokens.Scan() {
		return fmt.Errorf("extra output detected")
	}
	sp := shortestPath(c.n, edgesOut, c.s, c.t)
	if sp != c.L {
		return fmt.Errorf("path length %d != %d", sp, c.L)
	}
	if !expectPossible {
		// enumeration said impossible but candidate found one
		// accept it as valid
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		possible := existsAssignment(c)
		if err := runCase(bin, c, possible); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

const mod int64 = 1000000007

func scc(n int, g, gr [][]int) [][]int {
	visited := make([]bool, n)
	order := make([]int, 0, n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		visited[v] = true
		for _, u := range g[v] {
			if !visited[u] {
				dfs1(u)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}
	for i := range visited {
		visited[i] = false
	}
	comps := [][]int{}
	var dfs2 func(int, *[]int)
	dfs2 = func(v int, comp *[]int) {
		visited[v] = true
		*comp = append(*comp, v)
		for _, u := range gr[v] {
			if !visited[u] {
				dfs2(u, comp)
			}
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if !visited[v] {
			comp := []int{}
			dfs2(v, &comp)
			comps = append(comps, comp)
		}
	}
	return comps
}

func expected(n int, costs []int64, g, gr [][]int) (int64, int64) {
	comps := scc(n, g, gr)
	total := int64(0)
	ways := int64(1)
	for _, comp := range comps {
		minCost := costs[comp[0]]
		count := int64(0)
		for _, v := range comp {
			if costs[v] < minCost {
				minCost = costs[v]
				count = 1
			} else if costs[v] == minCost {
				count++
			}
		}
		total += minCost
		ways = (ways * count) % mod
	}
	return total, ways
}

func generateCase(rng *rand.Rand) (string, int64, int64) {
	n := rng.Intn(8) + 1
	m := rng.Intn(n*n + 1)
	costs := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		costs[i] = int64(rng.Intn(100) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", costs[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	g := make([][]int, n)
	gr := make([][]int, n)
	for i := 0; i < m; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		g[u] = append(g[u], v)
		gr[v] = append(gr[v], u)
		sb.WriteString(fmt.Sprintf("%d %d\n", u+1, v+1))
	}
	total, ways := expected(n, costs, g, gr)
	return sb.String(), total, ways
}

func runCase(exe, input string, expTotal, expWays int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var gotTotal, gotWays int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &gotTotal, &gotWays); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if gotTotal != expTotal || gotWays != expWays {
		return fmt.Errorf("expected %d %d got %d %d", expTotal, expWays, gotTotal, gotWays)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, t, w := generateCase(rng)
		if err := runCase(exe, in, t, w); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

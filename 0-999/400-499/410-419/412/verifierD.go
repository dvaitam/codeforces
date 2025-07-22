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

func expectedD(n int, edges [][2]int) string {
	g := make([][]int, n)
	for _, e := range edges {
		a, b := e[0]-1, e[1]-1
		g[a] = append(g[a], b)
	}
	vis := make([]bool, n)
	var res []int
	var dfs func(int)
	dfs = func(v int) {
		vis[v] = true
		for _, u := range g[v] {
			if !vis[u] {
				dfs(u)
			}
		}
		res = append(res, v+1)
	}
	for v := 0; v < n; v++ {
		if !vis[v] {
			dfs(v)
		}
	}
	strs := make([]string, len(res))
	for i, x := range res {
		strs[i] = fmt.Sprint(x)
	}
	return strings.Join(strs, " ")
}

func runCase(bin string, n int, edges [][2]int) error {
	input := fmt.Sprintf("%d %d\n", n, len(edges))
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e[0], e[1])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(expectedD(n, edges))
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(12) + 2
	maxE := n * (n - 1) / 2
	m := rng.Intn(maxE + 1)
	if m > 40 {
		m = 40
	}
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if used[[2]int{a, b}] || used[[2]int{b, a}] {
			continue
		}
		used[[2]int{a, b}] = true
		edges = append(edges, [2]int{a, b})
	}
	return n, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, edges := genCase(rng)
		if err := runCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

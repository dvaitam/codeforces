package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	indeg := make([]int, n+1)
	outdeg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		indeg[v]++
		outdeg[u]++
	}
	indeg2 := make([]int, n+1)
	copy(indeg2, indeg)
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if indeg2[i] == 0 {
			queue = append(queue, i)
		}
	}
	order := make([]int, 0, n)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, to := range adj[v] {
			indeg2[to]--
			if indeg2[to] == 0 {
				queue = append(queue, to)
			}
		}
	}
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = 1
	}
	ans := 1
	for _, v := range order {
		for _, to := range adj[v] {
			if outdeg[v] > 1 && indeg[to] > 1 {
				if dp[v]+1 > dp[to] {
					dp[to] = dp[v] + 1
					if dp[to] > ans {
						ans = dp[to]
					}
				}
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(48))
	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		edgeSet := make(map[[2]int]bool)
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			// ensure acyclicity by only adding if u < v randomly
			if rng.Intn(2) == 0 {
				if u > v {
					u, v = v, u
				}
			} else {
				if v > u {
					u, v = v, u
				}
			}
			e := [2]int{u, v}
			if !edgeSet[e] {
				edgeSet[e] = true
				edges = append(edges, e)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		exp := expected(n, edges)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

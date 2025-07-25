package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n int, edges [][2]int, specials []int) (int, int) {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	special := make([]bool, n+1)
	for _, s := range specials {
		special[s] = true
	}
	cnt := 0
	var dfs func(u, p int) bool
	dfs = func(u, p int) bool {
		need := special[u]
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			if dfs(v, u) {
				cnt++
				need = true
			}
		}
		return need
	}
	dfs(specials[0], 0)
	bfs := func(start int) (int, []int) {
		dist := make([]int, n+1)
		for i := 1; i <= n; i++ {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		far := start
		maxd := 0
		for head := 0; head < len(q); head++ {
			u := q[head]
			if special[u] {
				if dist[u] > maxd || (dist[u] == maxd && u < far) {
					far = u
					maxd = dist[u]
				}
			}
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		return far, dist
	}
	far1, _ := bfs(specials[0])
	far2, dist2 := bfs(far1)
	d := dist2[far2]
	startCity := far1
	if far2 < startCity {
		startCity = far2
	}
	time := 2*cnt - d
	return startCity, time
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	m := rng.Intn(n) + 1
	perm := rng.Perm(n)
	specials := make([]int, m)
	for i := 0; i < m; i++ {
		specials[i] = perm[i] + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i, s := range specials {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(s))
	}
	sb.WriteByte('\n')
	start, tm := expected(n, edges, specials)
	out := fmt.Sprintf("%d\n%d", start, tm)
	return sb.String(), out
}

func runCase(exe, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

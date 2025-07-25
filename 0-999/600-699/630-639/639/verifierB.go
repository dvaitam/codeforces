package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isPossible(n, d, h int) bool {
	if d < h {
		return false
	}
	if d > 2*h {
		return false
	}
	if d == 1 && n > 2 {
		return false
	}
	return true
}

func parseEdges(out string) ([]edge, error) {
	scan := bufio.NewScanner(strings.NewReader(out))
	edges := make([]edge, 0)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		var a, b int
		if _, err := fmt.Sscanf(line, "%d %d", &a, &b); err != nil {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		edges = append(edges, edge{a, b})
	}
	return edges, nil
}

func checkTree(n, d, h int, edges []edge) bool {
	if len(edges) != n-1 {
		return false
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		if e.u < 1 || e.u > n || e.v < 1 || e.v > n {
			return false
		}
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	vis := make([]bool, n+1)
	order := make([]int, 0, n)
	queue := []int{1}
	vis[1] = true
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, to := range adj[v] {
			if !vis[to] {
				vis[to] = true
				queue = append(queue, to)
			}
		}
	}
	if len(order) != n {
		return false
	}
	// compute height from 1
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	queue = []int{1}
	dist[1] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue = append(queue, to)
			}
		}
	}
	maxh := 0
	far := 1
	for i := 1; i <= n; i++ {
		if dist[i] > maxh {
			maxh = dist[i]
			far = i
		}
	}
	if maxh != h {
		return false
	}
	// diameter
	dist2 := make([]int, n+1)
	for i := range dist2 {
		dist2[i] = -1
	}
	queue = []int{far}
	dist2[far] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range adj[v] {
			if dist2[to] == -1 {
				dist2[to] = dist2[v] + 1
				queue = append(queue, to)
			}
		}
	}
	maxd := 0
	for i := 1; i <= n; i++ {
		if dist2[i] > maxd {
			maxd = dist2[i]
		}
	}
	return maxd == d
}

func generateCase(rng *rand.Rand) (string, bool) {
	n := rng.Intn(8) + 2
	d := rng.Intn(n-1) + 1
	h := rng.Intn(d) + 1
	input := fmt.Sprintf("%d %d %d\n", n, d, h)
	return input, isPossible(n, d, h)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, possible := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if !possible {
			if strings.TrimSpace(out) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:%s", i+1, out, input)
				os.Exit(1)
			}
			continue
		}
		var n, d, h int
		fmt.Sscanf(strings.TrimSpace(input), "%d %d %d", &n, &d, &h)
		edges, err2 := parseEdges(out)
		if err2 != nil || !checkTree(n, d, h, edges) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid tree\ninput:%s output:%s", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

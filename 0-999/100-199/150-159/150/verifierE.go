package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Edge struct{ to, w int }

func bestMedian(n, llim, rlim int, edges [][3]int) (int, [][2]int) {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		a, b, c := e[0], e[1], e[2]
		adj[a] = append(adj[a], Edge{b, c})
		adj[b] = append(adj[b], Edge{a, c})
	}
	best := -1
	pairs := [][2]int{}
	var dfs func(int, int, []int, []int)
	dfs = func(u, parent int, path []int, nodes []int) {
		nodes = append(nodes, u)
		if len(nodes) >= 2 {
			if len(nodes)-1 >= llim && len(nodes)-1 <= rlim {
				vals := append([]int(nil), path...)
				sort.Ints(vals)
				med := vals[len(vals)/2]
				if med > best {
					best = med
					pairs = pairs[:0]
					pairs = append(pairs, [2]int{nodes[0], u})
				} else if med == best {
					pairs = append(pairs, [2]int{nodes[0], u})
				}
			}
			if len(nodes)-1 == rlim {
				nodes = nodes[:len(nodes)-1]
				return
			}
		}
		for _, e := range adj[u] {
			if e.to == parent {
				continue
			}
			dfs(e.to, u, append(path, e.w), nodes)
		}
		nodes = nodes[:len(nodes)-1]
	}
	for i := 1; i <= n; i++ {
		dfs(i, -1, nil, nil)
	}
	if len(pairs) == 0 {
		pairs = append(pairs, [2]int{1, 1})
	}
	return best, pairs
}

func pairMedian(n int, edges [][3]int, x, y int) (int, int) {
	if x == y {
		return 0, 0
	}
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		a, b, c := e[0], e[1], e[2]
		adj[a] = append(adj[a], Edge{b, c})
		adj[b] = append(adj[b], Edge{a, c})
	}
	var path []int
	visited := make([]bool, n+1)
	var found bool
	var dfs func(int, int)
	dfs = func(u, parent int) {
		if found {
			return
		}
		if u == y {
			found = true
			return
		}
		visited[u] = true
		for _, e := range adj[u] {
			if e.to == parent || visited[e.to] {
				continue
			}
			path = append(path, e.w)
			dfs(e.to, u)
			if found {
				return
			}
			path = path[:len(path)-1]
		}
	}
	dfs(x, -1)
	if !found {
		return -1, 0
	}
	vals := append([]int(nil), path...)
	sort.Ints(vals)
	return vals[len(vals)/2], len(path)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		llim := rng.Intn(n-1) + 1
		rlim := rng.Intn(n-llim) + llim
		edges := make([][3]int, n-1)
		for j := 2; j <= n; j++ {
			p := rng.Intn(j-1) + 1
			w := rng.Intn(20)
			edges[j-2] = [3]int{p, j, w}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, llim, rlim))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		input := sb.String()
		bestMed, _ := bestMedian(n, llim, rlim, edges)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var x, y int
		if _, err := fmt.Sscanf(got, "%d %d", &x, &y); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output %q\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
		med, length := pairMedian(n, edges, x, y)
		if length < llim || length > rlim || med != bestMed {
			fmt.Fprintf(os.Stderr, "case %d failed: incorrect pair %d %d\ninput:\n%s", i+1, x, y, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (first number is count).
const embeddedTestcases = `100
5 2 1 1 2 3
4 6 3 2 3 4 1 3 1 2 2 3 4 1
4 6 1 2 4 1 3 3 2 4 2 3 3 1
4 2 1 3 1 4
3 4 2 2 1 1 3 3 3 2
5 4 5 3 2 4 3 3 4 5
4 4 4 2 1 1 3 2 2 3
5 4 2 3 3 4 4 5 1 1
2 3 2 1 2 2 1 1
5 4 2 1 4 2 3 2 4 5
3 2 3 2 2 1
4 0 
4 3 4 4 2 1 3 3
2 3 2 2 2 1 1 2
5 3 5 4 1 4 4 4
4 3 4 1 2 4 1 3
2 4 2 2 1 1 1 2 2 1
2 0 
4 4 4 2 3 1 3 2 1 2
3 4 1 3 3 3 1 1 1 2
3 3 1 1 3 1 3 3
3 3 2 1 2 2 1 1
3 4 2 2 1 1 2 1 3 1
2 1 2 1
4 5 1 1 3 2 2 2 3 3 4 2
2 4 1 2 2 1 2 2 1 1
3 3 3 1 1 2 2 1
4 5 2 3 3 4 4 4 3 3 1 4
2 4 1 2 2 2 2 1 1 1
2 4 1 2 1 1 2 2 2 1
3 5 2 1 1 2 2 2 2 3 1 3
3 1 3 2
3 1 3 2
5 1 5 3
5 4 2 2 2 4 1 1 3 3
5 5 2 3 5 4 5 3 1 1 5 1
5 1 2 1
3 1 3 2
5 4 4 1 3 2 5 3 3 5
3 6 2 2 3 1 3 2 3 3 2 3 2 1
4 3 3 1 2 3 1 4
3 0 
4 3 4 1 1 3 1 1
2 4 1 2 1 1 2 2 2 1
3 0 
3 6 3 3 2 3 1 1 1 2 2 1 2 2
5 5 3 3 2 2 3 5 5 5 5 4
2 0 
3 2 2 1 3 2
2 4 1 1 1 2 2 2 2 1
4 5 2 3 4 1 3 2 2 1 1 4
3 3 3 3 1 1 1 3
4 0 
5 2 4 2 4 3
2 1 1 1
5 5 3 5 2 5 3 4 4 3 3 1
2 3 1 1 2 1 2 2
3 0 
5 0 
4 0 
2 2 1 2 1 1
2 4 1 2 2 2 1 1 2 1
4 4 2 1 1 1 2 2 2 3
5 2 1 2 3 5
4 5 1 1 1 3 1 4 4 2 3 1
2 1 1 1
2 1 2 1
5 1 2 2
2 2 2 1 1 2
2 3 1 1 2 2 2 1
5 5 2 1 4 5 4 4 3 5 1 1
5 2 1 2 5 1
3 6 2 3 3 2 2 1 3 3 1 2 1 3
2 1 2 2
3 2 3 3 3 2
2 0 
2 2 1 1 1 2
4 3 3 3 3 1 3 4
2 0 
4 6 1 4 3 1 2 3 1 2 2 1 3 3
2 4 1 2 2 2 2 1 1 1
2 3 2 1 2 2 1 1
2 4 2 2 2 1 1 2 1 1
3 5 2 2 2 3 1 1 3 1 2 1
3 4 2 3 1 1 3 3 2 1
2 2 1 1 2 2
4 5 2 4 3 2 4 3 3 3 2 3
3 6 3 2 2 1 1 3 2 3 1 1 1 2
4 6 4 4 4 3 3 4 4 2 3 2 1 1
2 0 
4 0 
4 6 3 1 4 3 3 2 1 1 2 3 1 4
2 4 1 2 2 2 1 1 2 1
3 3 3 1 3 2 1 1
5 1 1 3
5 2 2 3 4 3
4 5 3 3 1 3 1 4 1 2 4 2
5 4 5 2 4 2 4 5 5 5
4 3 1 1 2 2 4 3
4 2 1 4 3 2`

// Hopcroft-Karp implementation from 387D.go.
type HopcroftKarp struct {
	n     int
	adj   [][]int
	pairU []int
	pairV []int
	dist  []int
	inf   int
}

func NewHK(n int, adj [][]int) *HopcroftKarp {
	return &HopcroftKarp{
		n:     n,
		adj:   adj,
		pairU: make([]int, n+1),
		pairV: make([]int, n+1),
		dist:  make([]int, n+1),
		inf:   1e9,
	}
}

func (hk *HopcroftKarp) bfs(center int) bool {
	queue := make([]int, 0, hk.n)
	for u := 1; u <= hk.n; u++ {
		if u == center {
			hk.dist[u] = hk.inf
		} else if hk.pairU[u] == 0 {
			hk.dist[u] = 0
			queue = append(queue, u)
		} else {
			hk.dist[u] = hk.inf
		}
	}
	hk.dist[0] = hk.inf
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		if hk.dist[u] < hk.dist[0] {
			for _, w := range hk.adj[u] {
				if w == center {
					continue
				}
				v := hk.pairV[w]
				if hk.dist[v] == hk.inf {
					hk.dist[v] = hk.dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
	}
	return hk.dist[0] != hk.inf
}

func (hk *HopcroftKarp) dfs(u, center int) bool {
	for _, w := range hk.adj[u] {
		if w == center {
			continue
		}
		v := hk.pairV[w]
		if hk.dist[v] == hk.dist[u]+1 {
			if v == 0 || hk.dfs(v, center) {
				hk.pairU[u] = w
				hk.pairV[w] = u
				return true
			}
		}
	}
	hk.dist[u] = hk.inf
	return false
}

func (hk *HopcroftKarp) MaxMatching(center int) int {
	for i := 1; i <= hk.n; i++ {
		hk.pairU[i] = 0
		hk.pairV[i] = 0
	}
	matching := 0
	for hk.bfs(center) {
		for u := 1; u <= hk.n; u++ {
			if u != center && hk.pairU[u] == 0 {
				if hk.dfs(u, center) {
					matching++
				}
			}
		}
	}
	return matching
}

func solveCase(n int, edges [][2]int) int {
	hasEdge := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		hasEdge[i] = make([]bool, n+1)
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		if !hasEdge[u][v] {
			hasEdge[u][v] = true
			adj[u] = append(adj[u], v)
		}
	}
	hk := NewHK(n, adj)
	m := len(edges)
	best := m + (3*n - 2)
	for center := 1; center <= n; center++ {
		cv := 0
		for u := 1; u <= n; u++ {
			if u == center {
				continue
			}
			if hasEdge[u][center] {
				cv++
			}
			if hasEdge[center][u] {
				cv++
			}
		}
		if hasEdge[center][center] {
			cv++
		}
		pv := hk.MaxMatching(center)
		k := cv + pv
		cost := m + (3*n - 2) - 2*k
		if cost < best {
			best = cost
		}
	}
	return best
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	fields := strings.Fields(embeddedTestcases)
	if len(fields) < 1 {
		fmt.Fprintln(os.Stderr, "no testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid testcase count: %v\n", err)
		os.Exit(1)
	}
	pos := 1
	for idx := 0; idx < t; idx++ {
		if pos+1 >= len(fields) {
			fmt.Fprintf(os.Stderr, "case %d: truncated header\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[pos])
		m, _ := strconv.Atoi(fields[pos+1])
		pos += 2
		if pos+2*m > len(fields) {
			fmt.Fprintf(os.Stderr, "case %d: missing edges\n", idx+1)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			pos += 2
			edges[i] = [2]int{u, v}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		expected := strconv.Itoa(solveCase(n, edges))
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	if pos != len(fields) {
		fmt.Fprintf(os.Stderr, "unused data after parsing tests\n")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}

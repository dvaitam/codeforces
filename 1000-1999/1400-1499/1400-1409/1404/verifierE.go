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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solve(reader *bufio.Reader) string {
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}
	hID := make([][]int, n)
	vID := make([][]int, n)
	for i := 0; i < n; i++ {
		hID[i] = make([]int, m)
		vID[i] = make([]int, m)
	}
	hCount := 0
	for i := 0; i < n; i++ {
		j := 0
		for j < m {
			if grid[i][j] == '#' {
				hCount++
				k := j
				for k < m && grid[i][k] == '#' {
					hID[i][k] = hCount
					k++
				}
				j = k
			} else {
				j++
			}
		}
	}
	vCount := 0
	for j := 0; j < m; j++ {
		i := 0
		for i < n {
			if grid[i][j] == '#' {
				vCount++
				k := i
				for k < n && grid[k][j] == '#' {
					vID[k][j] = vCount
					k++
				}
				i = k
			} else {
				i++
			}
		}
	}
	graph := make([][]int, hCount+1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				u := hID[i][j]
				v := vID[i][j]
				graph[u] = append(graph[u], v)
			}
		}
	}
	hk := NewHopcroftKarp(hCount, vCount, graph)
	ans := hk.MaxMatching()
	return fmt.Sprint(ans)
}

type HopcroftKarp struct {
	n, m  int
	graph [][]int
	pairU []int
	pairV []int
	dist  []int
}

func NewHopcroftKarp(n, m int, graph [][]int) *HopcroftKarp {
	return &HopcroftKarp{n: n, m: m, graph: graph, pairU: make([]int, n+1), pairV: make([]int, m+1), dist: make([]int, n+1)}
}

func (hk *HopcroftKarp) bfs() bool {
	const inf = 1 << 30
	queue := make([]int, 0, hk.n)
	for u := 1; u <= hk.n; u++ {
		if hk.pairU[u] == 0 {
			hk.dist[u] = 0
			queue = append(queue, u)
		} else {
			hk.dist[u] = inf
		}
	}
	found := false
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, v := range hk.graph[u] {
			if hk.pairV[v] == 0 {
				found = true
			} else if hk.dist[hk.pairV[v]] == inf {
				hk.dist[hk.pairV[v]] = hk.dist[u] + 1
				queue = append(queue, hk.pairV[v])
			}
		}
	}
	return found
}

func (hk *HopcroftKarp) dfs(u int) bool {
	const inf = 1 << 30
	for _, v := range hk.graph[u] {
		if hk.pairV[v] == 0 || (hk.dist[hk.pairV[v]] == hk.dist[u]+1 && hk.dfs(hk.pairV[v])) {
			hk.pairU[u] = v
			hk.pairV[v] = u
			return true
		}
	}
	hk.dist[u] = inf
	return false
}

func (hk *HopcroftKarp) MaxMatching() int {
	result := 0
	for hk.bfs() {
		for u := 1; u <= hk.n; u++ {
			if hk.pairU[u] == 0 && hk.dfs(u) {
				result++
			}
		}
	}
	return result
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('.')
			} else {
				sb.WriteByte('#')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

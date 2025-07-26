package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Graph struct {
	n   int
	adj [][]int
}

func runSolution(bin, input string) (string, error) {
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

func copyAdj(a [][]int) [][]int {
	n := len(a)
	b := make([][]int, n)
	for i := 0; i < n; i++ {
		b[i] = make([]int, n)
		copy(b[i], a[i])
	}
	return b
}

func addVertex(g Graph, u, v int) Graph {
	n := g.n + 1
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
	}
	for i := 0; i < g.n; i++ {
		copy(adj[i][:g.n], g.adj[i])
	}
	w := g.n
	adj[u][w]++
	adj[w][u]++
	adj[v][w]++
	adj[w][v]++
	return Graph{n, adj}
}

func canonical(g Graph) string {
	n := g.n
	perm := make([]int, n)
	perm[0], perm[1] = 0, 1
	best := ""
	used := make([]bool, n)
	used[0], used[1] = true, true
	var rec func(int)
	var sb strings.Builder
	rec = func(pos int) {
		if pos == n {
			sb.Reset()
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					sb.WriteString(fmt.Sprintf("%d,", g.adj[perm[i]][perm[j]]))
				}
				sb.WriteByte(';')
			}
			str := sb.String()
			if best == "" || str < best {
				best = str
			}
			return
		}
		for i := 2; i < n; i++ {
			if !used[i] {
				used[i] = true
				perm[pos] = i
				rec(pos + 1)
				used[i] = false
			}
		}
	}
	rec(2)
	return best
}

func expand(curr map[string]Graph) map[string]Graph {
	next := make(map[string]Graph)
	for _, g := range curr {
		for u := 0; u < g.n; u++ {
			for v := u + 1; v < g.n; v++ {
				if g.adj[u][v] > 0 {
					ng := addVertex(g, u, v)
					key := canonical(ng)
					if _, ok := next[key]; !ok {
						next[key] = ng
					}
				}
			}
		}
	}
	return next
}

func minCut(g Graph) int {
	edges := [][2]int{}
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			for k := 0; k < g.adj[i][j]; k++ {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	best := len(edges)
	for mask := 0; mask < (1 << len(edges)); mask++ {
		if bits.OnesCount(uint(mask)) >= best {
			continue
		}
		adj := make([][]bool, g.n)
		for i := range adj {
			adj[i] = make([]bool, g.n)
		}
		for idx, e := range edges {
			if mask&(1<<idx) != 0 {
				continue
			}
			adj[e[0]][e[1]] = true
			adj[e[1]][e[0]] = true
		}
		visited := make([]bool, g.n)
		queue := []int{0}
		visited[0] = true
		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			for y := 0; y < g.n; y++ {
				if adj[x][y] && !visited[y] {
					visited[y] = true
					queue = append(queue, y)
				}
			}
		}
		if !visited[1] {
			best = bits.OnesCount(uint(mask))
		}
	}
	return best
}

func countWorlds(n, m int) int {
	initAdj := [][]int{{0, 1}, {1, 0}}
	curr := map[string]Graph{"init": {2, initAdj}}
	for step := 0; step < n; step++ {
		curr = expand(curr)
	}
	cnt := 0
	for _, g := range curr {
		if minCut(g) == m {
			cnt++
		}
	}
	return cnt % 1000000007
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	m := rng.Intn(n+2) + 1
	return fmt.Sprintf("%d %d\n", n, m)
}

func verifyD(input, output string) error {
	var n, m int
	if _, err := fmt.Sscan(strings.TrimSpace(input), &n, &m); err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	expected := countWorlds(n, m)
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(output), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"1 1\n", "2 1\n", "2 2\n"}
	for len(cases) < 100 {
		cases = append(cases, generateCaseD(rng))
	}
	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyD(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

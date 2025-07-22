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

func shortestPath(n int, adj [][]int, a, b int) []int {
	prev := make([]int, n+1)
	for i := range prev {
		prev[i] = -1
	}
	q := []int{a}
	prev[a] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			break
		}
		for _, u := range adj[v] {
			if prev[u] != -1 {
				continue
			}
			prev[u] = v
			q = append(q, u)
		}
	}
	if prev[b] == -1 {
		return nil
	}
	path := []int{b}
	for b != a {
		b = prev[b]
		path = append(path, b)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func countPairs(n int, edges [][2]int) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	paths := make(map[[2]int][]int)
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			paths[[2]int{a, b}] = shortestPath(n, adj, a, b)
		}
	}
	var count int64
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			p1 := paths[[2]int{a, b}]
			set1 := make(map[int]bool)
			for _, x := range p1 {
				set1[x] = true
			}
			for c := 1; c <= n; c++ {
				for d := c + 1; d <= n; d++ {
					p2 := paths[[2]int{c, d}]
					ok := true
					for _, x := range p2 {
						if set1[x] {
							ok = false
							break
						}
					}
					if ok {
						count++
					}
				}
			}
		}
	}
	return count
}

func buildCase(n int, rng *rand.Rand) string {
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		u := i + 2
		v := rng.Intn(i+1) + 1
		edges[i] = [2]int{u, v}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func parseEdges(n int, s string) [][2]int {
	parts := strings.Fields(s)
	idx := 0
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		u := toInt(parts[idx])
		v := toInt(parts[idx+1])
		idx += 2
		edges[i] = [2]int{u, v}
	}
	return edges
}

func toInt(str string) int {
	var x int
	fmt.Sscan(str, &x)
	return x
}

func runCase(bin string, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	lines := strings.Fields(input)
	n := toInt(lines[0])
	edges := parseEdges(n, strings.Join(lines[1:], " "))
	exp := countPairs(n, edges)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"2\n1 2\n",
		"3\n1 2\n2 3\n",
	}
	for len(cases) < 100 {
		n := rng.Intn(6) + 2
		cases = append(cases, buildCase(n, rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

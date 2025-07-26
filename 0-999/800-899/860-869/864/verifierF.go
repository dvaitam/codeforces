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

func computePath(n int, edges [][2]int, s, t int) []int {
	adj := make([][]int, n+1)
	rev := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		rev[y] = append(rev[y], x)
	}
	for i := 1; i <= n; i++ {
		sort.Ints(adj[i])
	}
	reachable := make([]bool, n+1)
	queue := []int{t}
	reachable[t] = true
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, u := range rev[v] {
			if !reachable[u] {
				reachable[u] = true
				queue = append(queue, u)
			}
		}
	}
	if !reachable[s] {
		return nil
	}
	visited := make([]bool, n+1)
	curr := s
	visited[curr] = true
	path := []int{curr}
	for curr != t {
		found := false
		for _, v := range adj[curr] {
			if reachable[v] {
				curr = v
				if visited[curr] {
					return nil
				}
				visited[curr] = true
				path = append(path, curr)
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return path
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	m := rng.Intn(n * (n - 1))
	edges := make([][2]int, 0, m)
	edgeSet := make(map[[2]int]bool)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		e := [2]int{a, b}
		if edgeSet[e] {
			continue
		}
		edgeSet[e] = true
		edges = append(edges, e)
	}
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	var out strings.Builder
	for i := 0; i < q; i++ {
		s := rng.Intn(n) + 1
		t := rng.Intn(n) + 1
		for t == s {
			t = rng.Intn(n) + 1
		}
		k := rng.Intn(n) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", s, t, k))
		path := computePath(n, edges, s, t)
		if path == nil || k > len(path) {
			out.WriteString("-1\n")
		} else {
			out.WriteString(fmt.Sprintf("%d\n", path[k-1]))
		}
	}
	return sb.String(), strings.TrimSpace(out.String())
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

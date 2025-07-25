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

type pair struct{ a, b int }

func bfs(src int, adj [][]int) []int {
	n := len(adj)
	const inf = int(1e9)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[src] = 0
	q := []int{src}
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range adj[u] {
			if dist[v] > dist[u]+1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(n int, edges []pair, s1, t1, l1, s2, t2, l2 int) int {
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e.a, e.b
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(i, adj)
	}
	if dist[s1][t1] > l1 || dist[s2][t2] > l2 {
		return -1
	}
	ans := dist[s1][t1] + dist[s2][t2]
	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			d := dist[u][v]
			d1 := dist[s1][u] + d + dist[v][t1]
			if d1 > l1 {
				continue
			}
			d2 := dist[s2][u] + d + dist[v][t2]
			if d2 <= l2 {
				ans = min(ans, d1+d2-d)
			}
			d2 = dist[s2][v] + d + dist[u][t2]
			if d2 <= l2 {
				ans = min(ans, d1+d2-d)
			}
		}
	}
	m := len(edges)
	return m - ans
}

func genTestB() (string, int) {
	n := rand.Intn(5) + 2 //2..6
	// create random connected graph
	edges := make([]pair, 0)
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, pair{i, p})
	}
	extra := rand.Intn(n)
	edgeSet := map[pair]bool{}
	for _, e := range edges {
		if e.a > e.b {
			e.a, e.b = e.b, e.a
		}
		edgeSet[e] = true
	}
	for k := 0; k < extra; k++ {
		for {
			a := rand.Intn(n)
			b := rand.Intn(n)
			if a == b {
				continue
			}
			e := pair{a, b}
			if e.a > e.b {
				e.a, e.b = e.b, e.a
			}
			if edgeSet[e] {
				continue
			}
			edgeSet[e] = true
			edges = append(edges, pair{a, b})
			break
		}
	}
	s1 := rand.Intn(n)
	t1 := rand.Intn(n)
	s2 := rand.Intn(n)
	t2 := rand.Intn(n)
	l1 := rand.Intn(n) + 1
	l2 := rand.Intn(n) + 1
	input := fmt.Sprintf("%d %d\n", n, len(edges))
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e.a+1, e.b+1)
	}
	input += fmt.Sprintf("%d %d %d\n", s1+1, t1+1, l1)
	input += fmt.Sprintf("%d %d %d\n", s2+1, t2+1, l2)
	expected := solveB(n, edges, s1, t1, l1, s2, t2, l2)
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestB()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %d\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

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

type edge struct{ to int }

func makeTree(n int, rng *rand.Rand) ([][2]int, string) {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// build bracket sequence via DFS
	var sb strings.Builder
	var dfs func(int, int)
	dfs = func(u, p int) {
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			sb.WriteByte('(')
			dfs(v, u)
			sb.WriteByte(')')
		}
	}
	dfs(1, 0)
	return edges, sb.String()
}

func isValid(s []byte) bool {
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func buildTreeFromSeq(s []byte) [][]int {
	n := len(s)/2 + 1
	adj := make([][]int, n)
	stack := []int{0}
	id := 1
	for _, ch := range s {
		if ch == '(' {
			v := stack[len(stack)-1]
			u := id
			id++
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			stack = append(stack, u)
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return adj
}

func bfsFar(adj [][]int, start int) (int, int) {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	far := start
	for i, d := range dist {
		if d > dist[far] {
			far = i
		}
	}
	return far, dist[far]
}

func diameter(adj [][]int) int {
	a, _ := bfsFar(adj, 0)
	_, d := bfsFar(adj, a)
	return d
}

func genCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 3
	q := rng.Intn(4) + 1
	_, seq := makeTree(n, rng)
	bs := []byte(seq)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	sb.WriteString(seq)
	sb.WriteByte('\n')
	var out strings.Builder
	for i := 0; i < q; i++ {
		// attempt to find swap making valid sequence
		for {
			a := rng.Intn(len(bs))
			b := rng.Intn(len(bs))
			if a == b {
				continue
			}
			bs[a], bs[b] = bs[b], bs[a]
			if isValid(bs) {
				// record swap
				fmt.Fprintf(&sb, "%d %d\n", a+1, b+1)
				adj := buildTreeFromSeq(bs)
				out.WriteString(fmt.Sprintf("%d\n", diameter(adj)))
				break
			}
			bs[a], bs[b] = bs[b], bs[a]
		}
	}
	return sb.String(), out.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in, expect := genCaseC(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

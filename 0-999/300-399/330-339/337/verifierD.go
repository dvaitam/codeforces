package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func bfs(n int, g [][]int, start int) []int {
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, start)
	dist[start] = 0
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range g[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

func solve(n, m, d int, p []int, edges [][2]int) int {
	g := make([][]int, n)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	d0 := bfs(n, g, p[0])
	a := p[0]
	for _, x := range p {
		if d0[x] > d0[a] {
			a = x
		}
	}
	da := bfs(n, g, a)
	b := p[0]
	for _, x := range p {
		if da[x] > da[b] {
			b = x
		}
	}
	db := bfs(n, g, b)
	cnt := 0
	for i := 0; i < n; i++ {
		if da[i] <= d && db[i] <= d {
			cnt++
		}
	}
	return cnt
}

func generateTest(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	m := rng.Intn(n) + 1
	d := rng.Intn(n)
	p := make([]int, m)
	perm := rng.Perm(n)
	for i := 0; i < m; i++ {
		p[i] = perm[i]
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		parent := rng.Intn(i)
		edges[i-1] = [2]int{parent, i}
	}
	inp := fmt.Sprintf("%d %d %d\n", n, m, d)
	for i, x := range p {
		if i > 0 {
			inp += " "
		}
		inp += fmt.Sprintf("%d", x+1)
	}
	inp += "\n"
	for _, e := range edges {
		inp += fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1)
	}
	ans := solve(n, m, d, p, edges)
	return inp, ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inp)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n%s", t, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: failed to parse output: %v\nOutput: %s\n", t, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %d\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

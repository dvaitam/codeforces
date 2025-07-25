package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(n, m int, edges [][2]int) (int, int64) {
	g := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	color := make([]int8, n)
	for i := range color {
		color[i] = -1
	}
	var t int
	var ways int64
	var bipartitePairs int64
	queue := make([]int, 0, n)
	for v := 0; v < n; v++ {
		if color[v] != -1 {
			continue
		}
		queue = queue[:0]
		queue = append(queue, v)
		color[v] = 0
		cnt := [2]int64{1, 0}
		isBip := true
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, to := range g[u] {
				if color[to] == -1 {
					color[to] = color[u] ^ 1
					cnt[color[to]]++
					queue = append(queue, to)
				} else if color[to] == color[u] {
					isBip = false
				}
			}
		}
		if !isBip {
			return 0, 1
		}
		bipartitePairs += cnt[0]*(cnt[0]-1)/2 + cnt[1]*(cnt[1]-1)/2
	}
	if bipartitePairs > 0 {
		t = 1
		ways = bipartitePairs
	} else if m > 0 {
		t = 2
		ways = int64(m) * int64(n-2)
	} else {
		t = 3
		ways = int64(n) * int64(n-1) * int64(n-2) / 6
	}
	return t, ways
}

func genCase() (string, int, int64) {
	n := rand.Intn(6) + 1
	maxEdges := n * (n - 1) / 2
	m := rand.Intn(maxEdges + 1)
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		e := [2]int{a, b}
		if used[e] {
			continue
		}
		used[e] = true
		edges = append(edges, e)
	}
	t, ways := solveD(n, len(edges), edges)
	input := fmt.Sprintf("%d %d\n", n, len(edges))
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1)
	}
	return input, t, ways
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, et, ew := genCase()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var gt int
		var gw int64
		if _, err := fmt.Fscan(strings.NewReader(out.String()), &gt, &gw); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output\n", i+1)
			os.Exit(1)
		}
		if gt != et || gw != ew {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %d %d\ngot: %d %d\n", i+1, input, et, ew, gt, gw)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pair struct{ a, b int }

func mstAndMatching(n int, c int, w [][]int) int {
	const inf = int(1e9)
	inMST := make([]bool, n)
	minEdge := make([]int, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		minEdge[i] = inf
		parent[i] = -1
	}
	minEdge[0] = 0
	mstWeight := 0
	adj := make([][]int, n)
	for it := 0; it < n; it++ {
		u := -1
		for i := 0; i < n; i++ {
			if !inMST[i] && (u == -1 || minEdge[i] < minEdge[u]) {
				u = i
			}
		}
		inMST[u] = true
		if parent[u] != -1 {
			mstWeight += w[u][parent[u]]
			adj[u] = append(adj[u], parent[u])
			adj[parent[u]] = append(adj[parent[u]], u)
		}
		for v := 0; v < n; v++ {
			if !inMST[v] && w[u][v] > 0 && w[u][v] < minEdge[v] {
				minEdge[v] = w[u][v]
				parent[v] = u
			}
		}
	}

	var dfs func(int, int) (int, int)
	dfs = func(v, p int) (int, int) {
		children := []pair{}
		for _, to := range adj[v] {
			if to != p {
				c0, c1 := dfs(to, v)
				children = append(children, pair{c0, c1})
			}
		}
		dp1 := 0
		for _, ch := range children {
			if ch.a > ch.b {
				dp1 += ch.a
			} else {
				dp1 += ch.b
			}
		}
		dp0 := dp1
		for _, ch := range children {
			cand := dp1 - max(ch.a, ch.b) + ch.a + 1
			if cand > dp0 {
				dp0 = cand
			}
		}
		return dp0, dp1
	}
	m0, _ := dfs(0, -1)
	return mstWeight + m0*c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateCase() (int, int, [][]int) {
	n := rand.Intn(5) + 2 // 2..6
	c := rand.Intn(10) + 1
	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			val := rand.Intn(10) + 1
			w[i][j] = val
			w[j][i] = val
		}
	}
	return n, c, w
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	cs := make([]int, t)
	mats := make([][][]int, t)
	for i := 0; i < t; i++ {
		ns[i], cs[i], mats[i] = generateCase()
	}

	for idx := 0; idx < t; idx++ {
		n := ns[idx]
		c := cs[idx]
		w := mats[idx]
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, c)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprintf(&input, "%d", w[i][j])
			}
			input.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Println("failed to run binary:", err)
			os.Exit(1)
		}
		expected := mstAndMatching(n, c, w)
		var got int
		fmt.Sscan(strings.TrimSpace(out.String()), &got)
		if got != expected {
			fmt.Printf("case %d: expected %d, got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

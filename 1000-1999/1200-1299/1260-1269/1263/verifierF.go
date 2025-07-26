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

type edge struct {
	l, r int
}

func run(bin, input string) (string, error) {
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

func solveF(n int, parA []int, x []int, parB []int, y []int) int {
	a := len(parA) - 1
	b := len(parB) - 1
	chA := make([][]int, a+1)
	for i := 2; i <= a; i++ {
		p := parA[i]
		chA[p] = append(chA[p], i)
	}
	chB := make([][]int, b+1)
	for i := 2; i <= b; i++ {
		p := parB[i]
		chB[p] = append(chB[p], i)
	}
	leafIdxA := make([]int, a+1)
	leafIdxB := make([]int, b+1)
	for i := 1; i <= n; i++ {
		leafIdxA[x[i]] = i
		leafIdxB[y[i]] = i
	}
	lA := make([]int, a+1)
	rA := make([]int, a+1)
	var dfsA func(int)
	dfsA = func(u int) {
		if len(chA[u]) == 0 {
			idx := leafIdxA[u]
			lA[u], rA[u] = idx, idx
			return
		}
		l, r := n+1, 0
		for _, v := range chA[u] {
			dfsA(v)
			if lA[v] < l {
				l = lA[v]
			}
			if rA[v] > r {
				r = rA[v]
			}
		}
		lA[u], rA[u] = l, r
	}
	dfsA(1)
	lB := make([]int, b+1)
	rB := make([]int, b+1)
	var dfsB func(int)
	dfsB = func(u int) {
		if len(chB[u]) == 0 {
			idx := leafIdxB[u]
			lB[u], rB[u] = idx, idx
			return
		}
		l, r := n+1, 0
		for _, v := range chB[u] {
			dfsB(v)
			if lB[v] < l {
				l = lB[v]
			}
			if rB[v] > r {
				r = rB[v]
			}
		}
		lB[u], rB[u] = l, r
	}
	dfsB(1)
	edgesA := make([]edge, 0, a-1)
	for i := 2; i <= a; i++ {
		edgesA = append(edgesA, edge{lA[i], rA[i]})
	}
	edgesB := make([]edge, 0, b-1)
	for i := 2; i <= b; i++ {
		edgesB = append(edgesB, edge{lB[i], rB[i]})
	}
	n1 := len(edgesA)
	n2 := len(edgesB)
	adj := make([][]int, n1)
	for i := 0; i < n1; i++ {
		e1 := edgesA[i]
		for j := 0; j < n2; j++ {
			e2 := edgesB[j]
			if e1.l <= e2.r && e2.l <= e1.r {
				adj[i] = append(adj[i], j)
			}
		}
	}
	pairU := make([]int, n1)
	for i := range pairU {
		pairU[i] = -1
	}
	pairV := make([]int, n2)
	for i := range pairV {
		pairV[i] = -1
	}
	dist := make([]int, n1)
	INF := int(1e9)
	bfs := func() bool {
		q := []int{}
		for i := 0; i < n1; i++ {
			if pairU[i] == -1 {
				dist[i] = 0
				q = append(q, i)
			} else {
				dist[i] = INF
			}
		}
		found := false
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range adj[u] {
				pu := pairV[v]
				if pu != -1 && dist[pu] == INF {
					dist[pu] = dist[u] + 1
					q = append(q, pu)
				}
				if pu == -1 {
					found = true
				}
			}
		}
		return found
	}
	var dfs func(int) bool
	dfs = func(u int) bool {
		for _, v := range adj[u] {
			pu := pairV[v]
			if pu == -1 || (dist[pu] == dist[u]+1 && dfs(pu)) {
				pairU[u] = v
				pairV[v] = u
				return true
			}
		}
		dist[u] = INF
		return false
	}
	matching := 0
	for bfs() {
		for i := 0; i < n1; i++ {
			if pairU[i] == -1 {
				if dfs(i) {
					matching++
				}
			}
		}
	}
	totalEdges := (a - 1) + (b - 1)
	return totalEdges - matching
}

func generateTree(nLeaves int, rng *rand.Rand) (int, []int) {
	m := rng.Intn(2)
	a := 1 + m + nLeaves
	par := make([]int, a+1)
	for i := 2; i <= 1+m; i++ {
		par[i] = 1
	}
	for i := 0; i < nLeaves; i++ {
		id := m + 2 + i
		var parent int
		if i < 1+m {
			parent = 1 + i
			if parent > 1+m {
				parent = 1
			}
		} else {
			parent = rng.Intn(1+m) + 1
		}
		par[id] = parent
	}
	return a, par
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	a, parA := generateTree(n, rng)
	x := make([]int, n+1)
	for i := 1; i <= n; i++ {
		x[i] = len(parA) - n - 1 + i
	}
	b, parB := generateTree(n, rng)
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		y[i] = len(parB) - n - 1 + i
	}
	rng.Shuffle(n, func(i, j int) { y[i+1], y[j+1] = y[j+1], y[i+1] })
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	input.WriteString(fmt.Sprintf("%d\n", a))
	for i := 2; i <= a; i++ {
		input.WriteString(fmt.Sprintf("%d ", parA[i]))
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		input.WriteString(fmt.Sprintf("%d ", x[i]))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", b))
	for i := 2; i <= b; i++ {
		input.WriteString(fmt.Sprintf("%d ", parB[i]))
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		input.WriteString(fmt.Sprintf("%d ", y[i]))
	}
	input.WriteByte('\n')
	expect := fmt.Sprintf("%d", solveF(n, parA, x, parB, y))
	return input.String(), expect
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
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

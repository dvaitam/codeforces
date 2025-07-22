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

// Structures and globals similar to 294E.go

type edge struct{ to, w, id int }

var (
	adj       [][]edge
	edges     [][3]int
	visited   []bool
	compNodes []int
	sizeArr   []int
	downSum   []int64
	fullSum   []int64
)

func collect(u, skip int) {
	visited[u] = true
	compNodes = append(compNodes, u)
	for _, e := range adj[u] {
		if e.id == skip || visited[e.to] {
			continue
		}
		collect(e.to, skip)
	}
}

func dfs1(u, parent, skip int) {
	sizeArr[u] = 1
	downSum[u] = 0
	for _, e := range adj[u] {
		v := e.to
		if e.id == skip || v == parent || !visited[v] {
			continue
		}
		dfs1(v, u, skip)
		sizeArr[u] += sizeArr[v]
		downSum[u] += downSum[v] + int64(sizeArr[v])*int64(e.w)
	}
}

func dfs2(u, parent, compSize, skip int, minPtr *int64) {
	if fullSum[u] < *minPtr {
		*minPtr = fullSum[u]
	}
	for _, e := range adj[u] {
		v := e.to
		if e.id == skip || v == parent || !visited[v] {
			continue
		}
		fullSum[v] = fullSum[u] + int64(e.w)*int64(compSize-2*sizeArr[v])
		dfs2(v, u, compSize, skip, minPtr)
	}
}

func solveE(n int, edgesInput [][3]int) string {
	adj = make([][]edge, n)
	edges = edgesInput
	for i, e := range edges {
		a, b, w := e[0], e[1], e[2]
		adj[a] = append(adj[a], edge{b, w, i})
		adj[b] = append(adj[b], edge{a, w, i})
	}
	visited = make([]bool, n)
	sizeArr = make([]int, n)
	downSum = make([]int64, n)
	fullSum = make([]int64, n)
	for i := range visited {
		visited[i] = true
	}
	dfs1(0, -1, -1)
	fullSum[0] = downSum[0]
	minTmp := fullSum[0]
	dfs2(0, -1, n, -1, &minTmp)
	var sumAll int64
	for i := 0; i < n; i++ {
		sumAll += fullSum[i]
	}
	sumAll /= 2
	for i := range visited {
		visited[i] = false
	}
	sizeArr = make([]int, n)
	downSum = make([]int64, n)
	fullSum = make([]int64, n)
	bestDelta := int64(0)
	for id, e := range edges {
		u, v, _ := e[0], e[1], e[2]
		compNodes = compNodes[:0]
		collect(u, id)
		compSizeB := len(compNodes)
		dfs1(u, -1, id)
		fullSum[u] = downSum[u]
		minTmpB := fullSum[u]
		dfs2(u, -1, compSizeB, id, &minTmpB)
		sumBu := fullSum[u]
		for _, x := range compNodes {
			visited[x] = false
		}
		compNodes = compNodes[:0]
		collect(v, id)
		compSizeA := len(compNodes)
		dfs1(v, -1, id)
		fullSum[v] = downSum[v]
		minTmpA := fullSum[v]
		dfs2(v, -1, compSizeA, id, &minTmpA)
		sumAv := fullSum[v]
		for _, x := range compNodes {
			visited[x] = false
		}
		a := int64(compSizeA)
		b := int64(compSizeB)
		delta := b*(minTmpA-sumAv) + a*(minTmpB-sumBu)
		if id == 0 || delta < bestDelta {
			bestDelta = delta
		}
	}
	return fmt.Sprintf("%d\n", sumAll+bestDelta)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	edges := make([][3]int, n-1)
	parents := make([]int, n)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		parents[i] = p
		w := rng.Intn(10) + 1
		edges[i-1] = [3]int{p, i, w}
	}
	ans := solveE(n, edges)
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&in, "%d %d %d\n", e[0]+1, e[1]+1, e[2])
	}
	return in.String(), ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

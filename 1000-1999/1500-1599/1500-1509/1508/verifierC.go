package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded solver for 1508C: complement-graph BFS + Kruskal MST.
func solve1508C(input string) string {
	data := []byte(input)
	ptr := 0
	nextInt := func() int {
		for ptr < len(data) && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		val := 0
		for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
			val = val*10 + int(data[ptr]-'0')
			ptr++
		}
		return val
	}

	n := nextInt()
	m := nextInt()

	type Edge struct {
		u, v, w int
	}
	edges := make([]Edge, m)
	adj := make([]map[int]bool, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]bool)
	}
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := nextInt()
		edges[i] = Edge{u, v, w}
		adj[u][v] = true
		adj[v][u] = true
	}

	// Find connected components of complement graph using BFS with linked list
	next := make([]int, n+1)
	prev := make([]int, n+1)
	inList := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		prev[i] = i - 1
		next[i] = i + 1
		inList[i] = true
	}
	next[n] = 0
	head := 1

	remove := func(v int) {
		if !inList[v] {
			return
		}
		p := prev[v]
		nx := next[v]
		if p == 0 {
			head = nx
		} else {
			next[p] = nx
		}
		if nx != 0 {
			prev[nx] = p
		}
		inList[v] = false
	}

	comp := make([]int, n+1)
	compID := 0
	for s := 1; s <= n; s++ {
		if !inList[s] {
			continue
		}
		compID++
		remove(s)
		comp[s] = compID
		queue := []int{s}
		for qi := 0; qi < len(queue); qi++ {
			v := queue[qi]
			var toAdd []int
			for cur := head; cur != 0; {
				nx := next[cur]
				if !adj[v][cur] {
					remove(cur)
					comp[cur] = compID
					toAdd = append(toAdd, cur)
				}
				cur = nx
			}
			queue = append(queue, toAdd...)
		}
	}

	// DSU for Kruskal
	parent := make([]int, compID+1)
	rank := make([]int, compID+1)
	for i := 1; i <= compID; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) bool {
		a, b = find(a), find(b)
		if a == b {
			return false
		}
		if rank[a] < rank[b] {
			a, b = b, a
		}
		parent[b] = a
		if rank[a] == rank[b] {
			rank[a]++
		}
		return true
	}

	// Sort edges by weight, run Kruskal on component graph
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	totalWeight := int64(0)
	for _, e := range edges {
		cu, cv := comp[e.u], comp[e.v]
		if union(cu, cv) {
			totalWeight += int64(e.w)
		}
	}

	return strconv.FormatInt(totalWeight, 10)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges)
	if m == maxEdges {
		m--
	}
	edgesMap := make(map[[2]int]struct{})
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for len(edgesMap) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := edgesMap[key]; ok {
			continue
		}
		edgesMap[key] = struct{}{}
		w := r.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = bufio.NewWriter(os.Stdout)
	_ = io.Discard

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := solve1508C(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	to int
	w  int64
}

type Item struct {
	v    int
	dist int64
}

// Priority queue for Dijkstra
type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type test struct {
	input, expected string
}

func solve(input string) string {
	rdr := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(rdr, &n, &m, &k); err != nil {
		return ""
	}
	graph := make([][]Edge, n+1)
	roads := make([]struct {
		u, v int
		w    int64
	}, m)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(rdr, &u, &v, &w)
		graph[u] = append(graph[u], Edge{v, w})
		graph[v] = append(graph[v], Edge{u, w})
		roads[i] = struct {
			u, v int
			w    int64
		}{u, v, w}
	}
	trainsTo := make([][]int64, n+1)
	for i := 0; i < k; i++ {
		var s int
		var y int64
		fmt.Fscan(rdr, &s, &y)
		trainsTo[s] = append(trainsTo[s], y)
		graph[1] = append(graph[1], Edge{s, y})
		graph[s] = append(graph[s], Edge{1, y})
	}
	const INF int64 = 4e18
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	dist[1] = 0
	pq := &PQ{{1, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.v] {
			continue
		}
		for _, e := range graph[it.v] {
			nd := it.dist + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	hasRoad := make([]bool, n+1)
	for _, r := range roads {
		if dist[r.u]+r.w == dist[r.v] {
			hasRoad[r.v] = true
		}
		if dist[r.v]+r.w == dist[r.u] {
			hasRoad[r.u] = true
		}
	}
	var ans int64
	for v := 2; v <= n; v++ {
		var match int64
		for _, y := range trainsTo[v] {
			if y > dist[v] {
				ans++
			} else if y == dist[v] {
				match++
			}
		}
		if match > 0 {
			if hasRoad[v] {
				ans += match
			} else {
				ans += match - 1
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rand.Seed(450)
	var tests []test
	// fixed small cases
	fixed := []string{
		"2 1 1\n1 2 5\n2 3\n",
		"3 2 1\n1 2 1\n2 3 2\n3 2\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(6) + 2
		maxm := n * (n - 1) / 2
		m := rand.Intn(maxm) + 1
		k := rand.Intn(n) + 1
		edges := make([][3]int64, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			key1 := [2]int{u, v}
			key2 := [2]int{v, u}
			if used[key1] || used[key2] {
				continue
			}
			used[key1] = true
			used[key2] = true
			w := rand.Int63n(10) + 1
			edges = append(edges, [3]int64{int64(u), int64(v), w})
		}
		trains := make([][2]int64, k)
		for i := 0; i < k; i++ {
			s := rand.Intn(n-1) + 2
			y := rand.Int63n(10) + 1
			trains[i] = [2]int64{int64(s), y}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
		}
		for _, t := range trains {
			fmt.Fprintf(&sb, "%d %d\n", t[0], t[1])
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

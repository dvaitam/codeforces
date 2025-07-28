package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type Edge struct {
	to int
	w  int
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

const INF int64 = 1 << 60

func dijkstra(n int, revAdj [][]Edge, highlighted []bool) []int64 {
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	for i, hl := range highlighted {
		if hl {
			dist[i] = 0
			heap.Push(pq, Item{i, 0})
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		for _, e := range revAdj[u] {
			v := e.to
			nd := it.dist + int64(e.w)
			if nd < dist[v] {
				dist[v] = nd
				heap.Push(pq, Item{v, nd})
			}
		}
	}
	return dist
}

func computeCost(n int, adj, revAdj [][]Edge, highlighted []bool) int64 {
	any := false
	for _, hl := range highlighted {
		if hl {
			any = true
			break
		}
	}
	if !any {
		return -1
	}
	dist := dijkstra(n, revAdj, highlighted)
	cost := int64(0)
	for i := 0; i < n; i++ {
		if highlighted[i] {
			continue
		}
		if dist[i] == INF {
			return -1
		}
		best := int64(INF)
		for _, e := range adj[i] {
			if dist[e.to]+int64(e.w) == dist[i] {
				if int64(e.w) < best {
					best = int64(e.w)
				}
			}
		}
		if best == INF {
			return -1
		}
		cost += best
	}
	return cost
}

type testCaseH struct {
	n     int
	m     int
	q     int
	edges [][3]int
	ops   []struct {
		add bool
		v   int
	}
}

func genTestsH() []testCaseH {
	rng := rand.New(rand.NewSource(49))
	tests := make([]testCaseH, 100)
	for i := range tests {
		n := rng.Intn(4) + 3
		m := rng.Intn(5) + 1
		q := rng.Intn(4) + 1
		edges := make([][3]int, m)
		for j := 0; j < m; j++ {
			u := rng.Intn(n)
			v := rng.Intn(n)
			for v == u {
				v = rng.Intn(n)
			}
			w := rng.Intn(10) + 1
			edges[j] = [3]int{u, v, w}
		}
		ops := make([]struct {
			add bool
			v   int
		}, q)
		highlighted := make([]bool, n)
		for j := 0; j < q; j++ {
			if rng.Intn(2) == 0 {
				// add
				v := rng.Intn(n)
				for highlighted[v] {
					v = rng.Intn(n)
				}
				highlighted[v] = true
				ops[j] = struct {
					add bool
					v   int
				}{true, v}
			} else {
				// remove; ensure at least one highlighted
				v := rng.Intn(n)
				if !highlighted[v] {
					highlighted[v] = true
					ops[j] = struct {
						add bool
						v   int
					}{true, v}
				} else {
					ops[j] = struct {
						add bool
						v   int
					}{false, v}
					highlighted[v] = false
				}
			}
		}
		tests[i] = testCaseH{n, m, q, edges, ops}
	}
	return tests
}

func solveH(tc testCaseH) []int64 {
	adj := make([][]Edge, tc.n)
	rev := make([][]Edge, tc.n)
	for _, e := range tc.edges {
		u, v, w := e[0], e[1], e[2]
		adj[u] = append(adj[u], Edge{v, w})
		rev[v] = append(rev[v], Edge{u, w})
	}
	highlighted := make([]bool, tc.n)
	res := make([]int64, tc.q)
	for i, op := range tc.ops {
		if op.add {
			highlighted[op.v] = true
		} else {
			highlighted[op.v] = false
		}
		res[i] = computeCost(tc.n, adj, rev, highlighted)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsH()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, tc.q)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d %d\n", e[0]+1, e[1]+1, e[2])
		}
		for _, op := range tc.ops {
			if op.add {
				fmt.Fprintf(&input, "+ %d\n", op.v+1)
			} else {
				fmt.Fprintf(&input, "- %d\n", op.v+1)
			}
		}
		expected := solveH(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\noutput:\n%s\n", idx+1, err, out.String())
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
		scanner.Split(bufio.ScanWords)
		for j, exp := range expected {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "case %d wrong output format\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d non-integer output\n", idx+1)
				os.Exit(1)
			}
			if val != exp {
				fmt.Fprintf(os.Stderr, "case %d query %d: expected %d got %d\n", idx+1, j+1, exp, val)
				os.Exit(1)
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d extra output\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}

package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Edge struct{ to int }

type State struct {
	bad, dist int
	path      []int
	node      int
}

type PQ []*State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	a, b := pq[i], pq[j]
	if a.bad != b.bad {
		return a.bad < b.bad
	}
	if a.dist != b.dist {
		return a.dist < b.dist
	}
	pa, pb := a.path, b.path
	na, nb := len(pa), len(pb)
	lim := na
	if nb < lim {
		lim = nb
	}
	for i := 0; i < lim; i++ {
		if pa[i] != pb[i] {
			return pa[i] < pb[i]
		}
	}
	return na < nb
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PQ) Pop() interface{}   { old := *pq; n := len(old); x := old[n-1]; *pq = old[:n-1]; return x }

type testCase struct {
	input    string
	expected []string
}

const testcaseData = `100
7 4
6 3 5 1
4 6 7 2 3
5 3 1 2 1 5
4 5 4 7 1
3 1 4 6
5
+ 3 6
? 1 5
+ 1 3
? 6 3
? 1 3
4 3
1 4 3
3 1 2 4
4 4 3 2 3
5 3 2 2 2 1
4
+ 4 4
+ 1 4
+ 4 3
? 3 3
8 4
2 3 7 6
4 2 6 3 3
4 3 8 3 7
3 7 5 6
3 6 6 2
4
+ 6 7
+ 2 6
+ 3 3
+ 3 6
5 5
3 5 2 1 4
4 3 1 4 5
3 5 2 2
4 2 5 5 1
4 1 2 3 4
4 4 5 1 3
3
? 2 5
? 5 2
? 5 3
3 3
1 3 2
4 1 2 3 3
4 3 1 1 2
5 2 1 2 2 1
2
+ 2 2
? 2 3
3 3
2 3 1
3 2 2 3
4 3 1 2 1
5 1 1 3 2 2
4
? 3 3
+ 3 3
? 2 3
? 1 2
8 4
6 2 1 4
4 6 6 2 2
4 2 1 5 1
4 1 5 8 4
4 4 5 6 6
2
? 6 4
? 1 1
8 4
8 3 5 7
4 8 5 6 3
5 3 7 6 3 5
5 5 8 6 6 7
5 7 3 3 4 8
3
? 7 5
? 7 3
? 7 5
7 5
7 6 3 1 4
3 7 5 6
4 6 5 6 3
3 3 2 1
5 1 2 6 6 4
3 4 4 7
2
+ 6 7
+ 6 1
4 3
4 2 3
4 4 1 2 2
3 2 3 3
3 3 3 4
4
? 2 4
? 3 2
? 3 3
+ 2 3
7 4
5 3 6 2
4 5 4 7 3
5 3 6 5 3 6
4 6 2 2 2
5 2 4 5 4 5
2
? 5 6
? 5 2
3 3
3 1 2
5 3 3 2 2 1
4 1 1 2 2
4 2 2 3 3
4
+ 3 2
? 2 1
+ 1 3
+ 1 1
5 3
4 1 5
5 4 5 4 2 1
4 1 1 2 5
4 5 1 4 4
5
? 4 1
? 4 1
+ 4 5
? 1 1
? 1 5
8 5
2 1 8 3 5
5 2 6 5 1 1
4 1 6 6 8
5 8 3 1 1 3
4 3 8 2 5
5 5 4 1 7 2
4
? 3 3
+ 1 8
+ 2 8
? 3 5
5 5
4 3 1 5 2
5 4 5 4 3 3
5 3 2 4 5 1
5 1 1 1 2 5
4 5 1 1 2
5 2 1 3 4 4
1
? 5 4
3 3
1 3 2
4 1 3 2 3
3 3 1 2
5 2 1 2 1 1
1
? 1 2
3 3
3 1 2
3 3 3 1
5 1 1 3 2 2
5 2 1 1 2 3
1
+ 2 1
6 4
6 3 1 4
5 6 3 2 5 3
5 3 6 5 3 1
5 1 5 5 2 4
4 4 1 5 6
2
+ 1 3
? 1 4
8 4
3 2 6 1
5 3 6 2 1 2
3 2 2 6
5 6 8 4 7 1
4 1 8 6 3
1
? 3 1
3 3
2 1 3
3 2 3 1
4 1 1 1 3
5 3 3 2 2 2
2
+ 3 3
+ 1 1
4 4
3 4 2 1
5 3 1 2 2 4
4 4 1 1 2
3 2 1 1
3 1 2 3
2
+ 1 2
+ 1 2
7 4
3 4 6 1
3 3 2 4
5 4 2 1 3 6
4 6 7 4 1
3 1 5 3
3
+ 4 6
? 3 1
? 6 6
5 3
1 4 2
3 1 1 4
5 4 3 3 4 2
5 2 2 2 1 1
2
? 2 1
? 2 4
8 4
5 2 4 3
4 5 2 8 2
3 2 4 4
3 4 4 3
3 3 7 5
3
+ 5 3
? 3 2
+ 4 5
8 3
4 7 8
3 4 8 7
3 7 4 8
5 8 1 8 7 4
1
+ 7 7
5 5
1 5 4 3 2
3 1 5 5
3 5 1 4
4 4 1 1 3
5 3 2 5 1 2
5 2 2 5 1 1
4
? 3 3
? 3 2
? 1 4
? 5 2
3 3
2 3 1
4 2 2 2 3
3 3 1 1
5 1 1 2 1 2
1
? 1 1
4 4
1 2 4 3
5 1 2 3 3 2
3 2 1 4
4 4 3 4 3
5 3 1 3 2 1
5
+ 2 1
? 3 4
+ 3 2
+ 4 2
+ 3 4
5 3
1 5 4
3 1 1 5
3 5 3 4
4 4 3 5 1
1
+ 5 1
5 3
3 5 1
3 3 5 5
5 5 2 1 3 1
5 1 5 2 2 3
4
? 5 5
+ 5 5
? 5 1
+ 5 3
8 3
8 4 5
3 8 2 4
3 4 7 5
5 5 1 4 1 8
1
? 8 4
8 4
7 4 2 6
4 7 4 5 4
5 4 6 4 6 2
4 2 4 1 6
5 6 8 7 6 7
2
+ 4 6
+ 7 4
6 3
6 3 2
4 6 4 6 3
4 3 5 5 2
3 2 2 6
4
? 6 2
? 2 6
? 6 3
? 2 3
7 4
5 4 7 6
4 5 7 7 4
5 4 6 5 6 7
3 7 5 6
3 6 7 5
4
+ 5 4
+ 5 5
+ 7 6
+ 5 7
8 3
2 1 7
5 2 5 6 4 1
4 1 2 2 7
3 7 6 2
5
+ 2 7
+ 2 7
? 2 7
? 1 1
+ 7 7
8 5
6 2 1 7 3
4 6 5 2 2
3 2 2 1
5 1 7 3 4 7
5 7 2 7 7 3
3 3 7 6
4
? 2 1
+ 6 3
? 1 2
+ 2 6
5 5
5 2 4 1 3
5 5 1 3 1 2
3 2 4 4
3 4 5 1
3 1 3 3
3 3 4 5
4
? 2 3
? 1 5
+ 2 5
+ 2 4
7 3
5 4 3
3 5 6 4
5 4 5 6 7 3
3 3 3 5
1
? 5 5
5 4
2 4 3 1
3 2 4 4
5 4 4 4 1 3
3 3 5 1
4 1 2 1 2
4
+ 3 2
+ 3 2
+ 1 3
? 4 3
6 3
3 4 1
5 3 1 3 4 4
4 4 1 5 1
4 1 1 2 3
4
? 4 1
? 3 4
? 1 4
+ 3 1
6 3
1 4 3
3 1 1 4
3 4 2 3
4 3 1 1 1
3
? 4 3
+ 1 1
? 4 1
6 5
1 3 4 5 6
3 1 2 3
4 3 5 4 4
5 4 6 6 5 5
5 5 5 1 2 6
3 6 4 1
4
? 5 5
+ 3 4
+ 3 5
? 1 4
4 4
4 1 2 3
4 4 3 3 1
3 1 1 2
3 2 4 3
3 3 1 4
4
+ 3 1
+ 4 3
? 3 1
+ 4 2
7 3
1 2 4
3 1 3 2
3 2 7 4
5 4 3 3 7 1
2
+ 4 4
? 2 1
3 3
2 3 1
4 2 3 3 3
4 3 1 3 1
3 1 2 2
1
? 1 3
4 3
4 3 1
3 4 1 3
3 3 4 1
3 1 2 4
3
+ 1 4
+ 3 3
? 1 3
6 4
1 6 3 4
5 1 3 6 5 6
5 6 4 4 3 3
3 3 5 4
5 4 1 6 4 1
1
? 4 1
7 4
7 5 3 4
3 7 4 5
4 5 1 4 3
3 3 5 4
5 4 2 5 4 7
1
+ 7 7
5 3
3 5 2
4 3 4 2 5
5 5 1 3 3 2
3 2 3 3
4
+ 3 2
+ 5 3
? 2 3
? 5 2
3 3
2 1 3
5 2 3 1 1 1
5 1 2 1 2 3
5 3 3 3 2 2
5
? 3 3
+ 2 3
? 2 2
+ 1 1
? 3 3
3 3
1 2 3
4 1 2 2 2
4 2 1 2 3
5 3 3 2 2 1
4
+ 1 2
+ 3 1
+ 1 2
? 3 2
6 5
5 4 2 1 3
4 5 2 6 4
5 4 4 5 3 2
3 2 2 1
4 1 3 5 3
4 3 5 5 5
2
? 4 1
? 2 2
5 5
1 2 4 3 5
5 1 3 4 1 2
3 2 1 4
4 4 1 1 3
5 3 3 5 2 5
3 5 5 1
4
+ 4 4
+ 4 1
+ 4 4
+ 3 1
3 3
1 3 2
3 1 3 3
5 3 2 1 2 2
4 2 1 1 1
4
+ 2 2
+ 1 2
? 3 1
+ 1 1
4 3
1 4 2
4 1 1 3 4
4 4 4 4 2
3 2 2 1
5
? 1 4
? 1 1
+ 2 2
+ 4 1
? 2 4
5 4
5 3 2 1
5 5 3 4 4 3
4 3 5 1 2
3 2 5 1
3 1 4 5
1
+ 3 3
4 3
1 3 2
4 1 4 2 3
5 3 4 1 3 2
5 2 3 3 1 1
5
? 1 1
? 3 2
? 2 3
? 2 2
+ 3 1
5 3
5 4 2
4 5 3 2 4
4 4 2 1 2
3 2 4 5
2
+ 4 2
+ 2 2
6 4
6 5 3 4
5 6 5 1 2 5
5 5 4 6 1 3
5 3 4 4 3 4
3 4 2 6
5
+ 6 5
+ 5 5
? 4 6
+ 5 3
? 5 5
7 3
6 5 3
4 6 5 3 5
3 5 6 3
3 3 3 6
5
? 5 3
? 5 3
? 6 5
+ 6 6
? 6 6
6 4
4 6 3 5
3 4 1 6
5 6 5 5 4 3
5 3 5 5 6 5
5 5 2 1 6 4
4
? 6 4
? 6 4
? 4 6
+ 4 4
6 3
3 2 4
4 3 2 1 2
4 2 6 2 4
3 4 2 3
4
? 3 2
? 2 2
? 3 4
? 4 2
3 3
2 3 1
5 2 1 1 2 3
5 3 2 3 3 1
3 1 1 2
5
? 2 3
? 1 1
+ 2 1
+ 2 1
+ 3 2
8 3
8 7 3
5 8 2 8 3 7
4 7 1 1 3
5 3 6 8 5 8
5
+ 7 3
? 7 8
+ 3 7
? 8 8
+ 8 3
3 3
1 3 2
3 1 1 3
5 3 2 2 1 2
4 2 2 1 1
3
? 1 2
? 2 2
? 1 3
7 4
2 4 3 7
4 2 3 5 4
4 4 5 1 3
5 3 5 2 3 7
3 7 5 2
3
+ 7 3
+ 4 3
? 7 2
5 3
3 1 2
4 3 1 4 1
4 1 1 4 2
5 2 1 5 2 3
3
+ 1 2
? 2 2
+ 2 2
6 4
6 2 4 5
3 6 2 2
3 2 4 4
4 4 1 3 5
4 5 4 5 6
5
+ 6 5
+ 4 6
? 6 6
? 4 4
+ 2 6
6 4
3 2 6 4
5 3 5 5 4 2
4 2 5 4 6
3 6 2 4
4 4 4 3 3
5
+ 4 6
+ 6 3
+ 4 6
+ 4 6
+ 2 4
3 3
2 3 1
5 2 1 2 3 3
5 3 2 1 2 1
4 1 1 1 2
2
? 2 1
+ 1 1
7 3
1 5 3
3 1 3 5
4 5 5 3 3
3 3 7 1
1
+ 5 1
5 4
5 2 1 3
4 5 3 2 2
4 2 4 1 1
3 1 1 3
3 3 1 5
3
+ 5 5
? 1 1
+ 2 1
5 3
4 2 1
5 4 4 1 3 2
4 2 2 1 1
5 1 5 2 4 4
5
+ 2 4
? 4 2
+ 1 4
? 2 2
+ 1 2
3 3
3 2 1
4 3 2 2 2
5 2 1 3 2 1
4 1 1 3 3
1
? 3 3
3 3
1 3 2
3 1 2 3
5 3 3 3 1 2
4 2 3 1 1
2
? 1 1
? 2 2
8 5
8 2 4 3 5
3 8 3 2
5 2 8 2 6 4
5 4 8 1 5 3
3 3 3 5
5 5 3 5 7 8
3
? 5 5
? 2 2
+ 5 2
8 3
1 4 2
4 1 1 4 4
4 4 3 5 2
4 2 5 2 1
3
+ 1 2
+ 2 1
? 4 4
5 3
4 5 2
3 4 4 5
5 5 4 5 1 2
4 2 5 5 4
1
+ 2 5
5 5
2 3 5 1 4
5 2 2 5 5 3
5 3 1 3 2 5
4 5 4 4 1
5 1 2 4 4 4
4 4 4 5 2
2
? 5 3
? 2 2
6 5
4 5 1 3 6
3 4 4 5
5 5 5 6 6 1
3 1 1 3
3 3 1 6
4 6 4 2 4
1
? 1 3
7 5
5 6 3 2 1
4 5 5 4 6
5 6 5 3 1 3
5 3 3 1 4 2
3 2 2 1
3 1 5 5
3
+ 3 6
? 6 2
? 1 2
6 3
3 4 6
3 3 3 4
4 4 1 3 6
3 6 2 3
2
+ 3 4
? 4 3
8 5
8 3 1 5 7
3 8 3 3
5 3 2 2 8 1
4 1 3 4 5
4 5 6 3 7
3 7 7 8
4
+ 8 1
? 7 5
? 8 8
+ 7 5
4 3
3 4 2
4 3 1 4 4
4 4 1 4 2
4 2 3 4 3
1
+ 3 3
6 3
1 5 2
5 1 5 5 4 5
5 5 2 6 6 2
5 2 4 5 5 1
3
+ 2 1
+ 2 5
+ 2 5
5 5
2 3 5 4 1
3 2 4 3
4 3 3 1 5
4 5 1 5 4
5 4 5 3 4 1
3 1 4 2
1
? 2 3
6 5
5 2 6 4 1
5 5 1 5 4 2
3 2 2 6
4 6 6 4 4
5 4 2 4 2 1
4 1 4 1 5
4
+ 4 6
? 6 1
? 2 5
? 6 6
8 5
1 6 3 7 5
3 1 4 6
4 6 7 6 3
5 3 5 7 7 7
3 7 3 5
5 5 1 6 1 1
1
? 7 5
7 4
6 7 1 4
3 6 3 7
3 7 4 1
3 1 5 4
3 4 4 6
3
? 1 7
+ 1 1
+ 7 7
3 3
2 1 3
5 2 2 2 2 1
3 1 2 3
3 3 3 2
3
? 1 3
+ 2 1
+ 1 3
7 5
1 6 5 2 7
4 1 7 2 6
5 6 1 6 6 5
5 5 3 6 3 2
3 2 7 7
3 7 7 1
5
+ 2 1
? 6 6
? 5 5
? 6 7
+ 6 6
7 4
3 5 4 6
4 3 6 3 5
4 5 3 1 4
4 4 6 4 6
4 6 5 5 3
2
+ 3 3
+ 3 6
8 5
2 7 3 8 4
4 2 1 8 7
5 7 6 6 5 3
3 3 4 8
5 8 2 8 3 4
5 4 2 6 5 2
2
? 4 7
? 8 2
4 4
1 4 2 3
4 1 1 1 4
5 4 2 2 4 2
5 2 4 3 1 3
3 3 2 1
5
+ 2 3
+ 4 3
? 4 1
+ 1 3
+ 2 3
8 4
8 4 2 3
4 8 8 7 4
3 4 1 2
5 2 1 4 8 3
3 3 6 8
3
? 3 8
+ 3 8
+ 8 8
8 5
3 1 2 4 5
4 3 7 4 1
5 1 6 4 4 2
3 2 5 4
4 4 1 7 5
5 5 7 4 6 3
2
? 1 3
+ 1 3
3 3
2 3 1
4 2 1 2 3
5 3 1 3 2 1
3 1 3 2
5
+ 1 2
? 1 3
? 2 2
+ 2 1
? 1 3
3 3
2 3 1
3 2 1 3
5 3 1 2 2 1
4 1 1 2 2
5
+ 2 2
? 2 2
+ 1 3
? 2 1
? 1 3
4 3
2 1 4
3 2 2 1
4 1 3 2 4
4 4 1 3 2
3
? 1 4
? 4 4
+ 1 4
4 4
2 4 1 3
5 2 3 3 2 4
3 4 4 1
3 1 1 3
5 3 2 3 3 2
3
+ 3 4
? 3 4
+ 4 1`

func solveCase(n, m int, good []int, shortcuts [][]int, queries []string) []string {
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := good[i]
		v := good[(i+1)%m]
		adj[u] = append(adj[u], Edge{v})
		adj[v] = append(adj[v], Edge{u})
	}
	for _, pts := range shortcuts {
		for j := 0; j+1 < len(pts); j++ {
			u, v := pts[j], pts[j+1]
			adj[u] = append(adj[u], Edge{v})
			adj[v] = append(adj[v], Edge{u})
		}
	}
	for i := 1; i <= n; i++ {
		sort.Slice(adj[i], func(a, b int) bool { return adj[i][a].to < adj[i][b].to })
	}
	weights := make(map[int]int)
	res := []string{}
	for _, line := range queries {
		parts := strings.Fields(line)
		op := parts[0]
		s, _ := strconv.Atoi(parts[1])
		t, _ := strconv.Atoi(parts[2])
		if op == "+" {
			key1 := s*(n+1) + t
			key2 := t*(n+1) + s
			weights[key1]++
			weights[key2]++
		} else {
			dist := make([]int, n+1)
			badc := make([]int, n+1)
			vis := make([]bool, n+1)
			paths := make([][]int, n+1)
			for i := 1; i <= n; i++ {
				dist[i] = 1e9
				badc[i] = 1e9
			}
			pq := &PQ{}
			heap.Init(pq)
			badc[s], dist[s] = 0, 0
			paths[s] = []int{s}
			heap.Push(pq, &State{0, 0, []int{s}, s})
			var ans *State
			for pq.Len() > 0 {
				cur := heap.Pop(pq).(*State)
				u := cur.node
				if vis[u] {
					continue
				}
				vis[u] = true
				if u == t {
					ans = cur
					break
				}
				for _, e := range adj[u] {
					v := e.to
					if vis[v] {
						continue
					}
					key := u*(n+1) + v
					w := weights[key]
					nb := cur.bad + w
					nd := cur.dist + 1
					newPath := append(append([]int{}, cur.path...), v)
					better := false
					if nb < badc[v] || (nb == badc[v] && nd < dist[v]) {
						better = true
					}
					if nb == badc[v] && nd == dist[v] {
						pOld := paths[v]
						for x := 0; x < len(newPath) && x < len(pOld); x++ {
							if newPath[x] != pOld[x] {
								if newPath[x] < pOld[x] {
									better = true
								}
								break
							}
						}
						if !better && len(newPath) < len(pOld) {
							better = true
						}
					}
					if better {
						badc[v], dist[v] = nb, nd
						paths[v] = newPath
						heap.Push(pq, &State{nb, nd, newPath, v})
					}
				}
			}
			if ans == nil {
				res = append(res, "-1")
			} else {
				res = append(res, fmt.Sprintf("%d", ans.bad))
				for j := 1; j < len(ans.path); j++ {
					u := ans.path[j-1]
					v := ans.path[j]
					key1 := u*(n+1) + v
					key2 := v*(n+1) + u
					delete(weights, key1)
					delete(weights, key2)
				}
			}
		}
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: incomplete header", caseNum+1)
		}
		n, _ := strconv.Atoi(fields[pos])
		m, _ := strconv.Atoi(fields[pos+1])
		pos += 2
		good := make([]int, m)
		for i := 0; i < m; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing good", caseNum+1)
			}
			v, _ := strconv.Atoi(fields[pos])
			pos++
			good[i] = v
		}
		shortcuts := make([][]int, m)
		for i := 0; i < m; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing shortcut len", caseNum+1)
			}
			k, _ := strconv.Atoi(fields[pos])
			pos++
			pts := make([]int, k)
			for j := 0; j < k; j++ {
				if pos >= len(fields) {
					return nil, fmt.Errorf("case %d: missing shortcut pt", caseNum+1)
				}
				val, _ := strconv.Atoi(fields[pos])
				pos++
				pts[j] = val
			}
			shortcuts[i] = pts
		}
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing query count", caseNum+1)
		}
		q, _ := strconv.Atoi(fields[pos])
		pos++
		queries := make([]string, q)
		for i := 0; i < q; i++ {
			if pos+2 >= len(fields) {
				return nil, fmt.Errorf("case %d: incomplete query", caseNum+1)
			}
			op := fields[pos]
			s := fields[pos+1]
			t := fields[pos+2]
			pos += 3
			queries[i] = fmt.Sprintf("%s %s %s", op, s, t)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range good {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			sb.WriteString(strconv.Itoa(len(shortcuts[i])))
			for _, v := range shortcuts[i] {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(v))
			}
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(q))
		sb.WriteByte('\n')
		for _, ql := range queries {
			sb.WriteString(ql)
			sb.WriteByte('\n')
		}
		expected := solveCase(n, m, good, shortcuts, queries)
		cases = append(cases, testCase{input: sb.String(), expected: expected})
	}
	return cases, nil
}

func run(bin, input string) ([]string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.Fields(strings.TrimSpace(out.String())), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := tc.expected
		if len(got) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d outputs got %d\n", idx+1, len(exp), len(got))
			os.Exit(1)
		}
		for i := range exp {
			if got[i] != exp[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, exp[i], got[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

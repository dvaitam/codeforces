package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	d  int
}

func digits(num int) []int {
	if num == 0 {
		return []int{0}
	}
	var ds [10]int
	n := 0
	for num > 0 {
		ds[n] = num % 10
		n++
		num /= 10
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = ds[n-1-i]
	}
	return res
}

const MOD int = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	// Build expanded graph
	// Preallocate rough size: n + m*6*2
	maxNodes := n + m*12
	adj := make([][]Edge, maxNodes+5)
	nextID := n
	for i := 1; i <= m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		ds := digits(i)
		// x -> y
		cur := x
		for j := 0; j < len(ds)-1; j++ {
			nextID++
			adj[cur] = append(adj[cur], Edge{to: nextID, d: ds[j]})
			cur = nextID
		}
		adj[cur] = append(adj[cur], Edge{to: y, d: ds[len(ds)-1]})
		// y -> x
		cur = y
		for j := 0; j < len(ds)-1; j++ {
			nextID++
			adj[cur] = append(adj[cur], Edge{to: nextID, d: ds[j]})
			cur = nextID
		}
		adj[cur] = append(adj[cur], Edge{to: x, d: ds[len(ds)-1]})
	}
	total := nextID

	// BFS for distance
	dist := make([]int, total+1)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, total)
	q = append(q, 1)
	dist[1] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, e := range adj[v] {
			if dist[e.to] == -1 {
				dist[e.to] = dist[v] + 1
				q = append(q, e.to)
			}
		}
	}

	// BFS for minimal lexicographic number
	val := make([]int, total+1)
	vis := make([]bool, total+1)
	curLayer := []int{1}
	vis[1] = true
	for len(curLayer) > 0 {
		buckets := [10]map[int]int{}
		for i := 0; i < 10; i++ {
			buckets[i] = make(map[int]int)
		}
		for _, v := range curLayer {
			for _, e := range adj[v] {
				if dist[e.to] == dist[v]+1 {
					cand := (val[v]*10 + e.d) % MOD
					m := buckets[e.d]
					if old, ok := m[e.to]; !ok || cand < old {
						m[e.to] = cand
					}
				}
			}
		}
		nextLayer := make([]int, 0)
		for d := 0; d <= 9; d++ {
			m := buckets[d]
			for v, c := range m {
				if !vis[v] {
					vis[v] = true
					val[v] = c
					nextLayer = append(nextLayer, v)
				}
			}
		}
		curLayer = nextLayer
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 2; i <= n; i++ {
		fmt.Fprintf(out, "%d\n", val[i])
	}
}

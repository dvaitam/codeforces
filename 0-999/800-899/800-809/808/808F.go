package main

import (
	"bufio"
	"fmt"
	"os"
)

// simple Dinic implementation with int64 capacities

type Edge struct {
	to  int
	rev int
	cap int64
}

type Dinic struct {
	g     [][]Edge
	level []int
	iter  []int
}

func NewDinic(n int) *Dinic {
	g := make([][]Edge, n)
	level := make([]int, n)
	iter := make([]int, n)
	return &Dinic{g: g, level: level, iter: iter}
}

func (d *Dinic) AddEdge(from, to int, cap int64) {
	d.g[from] = append(d.g[from], Edge{to, len(d.g[to]), cap})
	d.g[to] = append(d.g[to], Edge{from, len(d.g[from]) - 1, 0})
}

func (d *Dinic) bfs(s int, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0, len(d.g))
	d.level[s] = 0
	q = append(q, s)
	for h := 0; h < len(q); h++ {
		v := q[h]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				q = append(q, e.to)
				if e.to == t {
					return true
				}
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(v, t int, f int64) int64 {
	if v == t {
		return f
	}
	for ; d.iter[v] < len(d.g[v]); d.iter[v]++ {
		e := &d.g[v][d.iter[v]]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min64(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.g[e.to][e.rev].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
	var flow int64
	const INF int64 = 1 << 60
	for d.bfs(s, t) {
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, INF)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

const maxC = 200000

var isPrime [maxC + 1]bool

func sieve() {
	for i := 2; i <= maxC; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= maxC; i++ {
		if isPrime[i] {
			for j := i * i; j <= maxC; j += i {
				isPrime[j] = false
			}
		}
	}
}

type Card struct {
	p int
	c int
	l int
}

var n, k int
var cards []Card

func feasible(level int) bool {
	type item struct{ p, c int }
	nodes := make([]item, 0)
	max1 := -1
	for _, card := range cards {
		if card.l <= level {
			if card.c == 1 {
				if card.p > max1 {
					max1 = card.p
				}
			} else {
				nodes = append(nodes, item{card.p, card.c})
			}
		}
	}
	if max1 > 0 {
		nodes = append(nodes, item{max1, 1})
	}
	if len(nodes) == 0 {
		return false
	}
	total := 0
	oddIdx := []int{}
	evenIdx := []int{}
	for i, it := range nodes {
		total += it.p
		if it.c%2 == 1 {
			oddIdx = append(oddIdx, i)
		} else {
			evenIdx = append(evenIdx, i)
		}
	}
	m := len(nodes)
	S := m
	T := m + 1
	d := NewDinic(m + 2)
	for i, it := range nodes {
		if it.c%2 == 1 {
			d.AddEdge(S, i, int64(it.p))
		} else {
			d.AddEdge(i, T, int64(it.p))
		}
	}
	const INF int64 = 1 << 50
	for _, oi := range oddIdx {
		for _, ej := range evenIdx {
			if isPrime[nodes[oi].c+nodes[ej].c] {
				d.AddEdge(oi, ej, INF)
			}
		}
	}
	maxflow := d.MaxFlow(S, T)
	remaining := int64(total) - maxflow
	return remaining >= int64(k)
}

func main() {
	sieve()
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	cards = make([]Card, n)
	maxL := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &cards[i].p, &cards[i].c, &cards[i].l)
		if cards[i].l > maxL {
			maxL = cards[i].l
		}
	}
	if !feasible(maxL) {
		fmt.Println(-1)
		return
	}
	lo, hi := 1, maxL
	for lo < hi {
		mid := (lo + hi) / 2
		if feasible(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Println(lo)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Bitset represents a set of nodes using 64-bit words
// size must not exceed 256 bits (enough for this problem)
type Bitset []uint64

func NewBitset(n int) Bitset {
	return make([]uint64, (n+63)>>6)
}

func (b Bitset) Set(i int)      { b[i>>6] |= 1 << (uint(i) & 63) }
func (b Bitset) Has(i int) bool { return (b[i>>6]>>(uint(i)&63))&1 != 0 }
func (b Bitset) Or(other Bitset) {
	for i := range b {
		b[i] |= other[i]
	}
}
func (b Bitset) Clone() Bitset {
	c := make([]uint64, len(b))
	copy(c, b)
	return c
}

// Matrix is a boolean adjacency matrix using bitsets for rows
// m[i] is the set of nodes reachable from i in one step

type Matrix []Bitset

func NewMatrix(n int) Matrix {
	m := make(Matrix, n)
	for i := range m {
		m[i] = NewBitset(n)
	}
	return m
}

func (m Matrix) AddEdge(a, b int) { m[a].Set(b) }

func multiplyVecMat(v Bitset, m Matrix) Bitset {
	res := NewBitset(len(m))
	for i := 0; i < len(m); i++ {
		if v.Has(i) {
			for w := 0; w < len(res); w++ {
				res[w] |= m[i][w]
			}
		}
	}
	return res
}

func multiplyMatMat(a, b Matrix) Matrix {
	n := len(a)
	res := NewMatrix(n)
	for i := 0; i < n; i++ {
		row := a[i]
		for j := 0; j < n; j++ {
			if row.Has(j) {
				for w := 0; w < len(res[i]); w++ {
					res[i][w] |= b[j][w]
				}
			}
		}
	}
	return res
}

func buildPowers(base Matrix, log int) []Matrix {
	pow := make([]Matrix, log)
	pow[0] = base
	for i := 1; i < log; i++ {
		pow[i] = multiplyMatMat(pow[i-1], pow[i-1])
	}
	return pow
}

func applyPower(v Bitset, pow []Matrix, k int64) Bitset {
	res := v.Clone()
	idx := 0
	for k > 0 {
		if k&1 == 1 {
			res = multiplyVecMat(res, pow[idx])
		}
		k >>= 1
		idx++
	}
	return res
}

func earliestReach(v Bitset, pow []Matrix, limit int64, target int) (int64, bool) {
	if v.Has(target) {
		return 0, true
	}
	l, r := int64(1), limit
	ans := int64(-1)
	for l <= r {
		mid := (l + r) / 2
		tmp := applyPower(v, pow, mid)
		if tmp.Has(target) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	if ans == -1 {
		return 0, false
	}
	return ans, true
}

func bfsFromSet(start Bitset, adj [][]int, target int) int {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0)
	for i := 0; i < n; i++ {
		if start.Has(i) {
			dist[i] = 0
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		v := q[head]
		if v == target {
			return dist[v]
		}
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	type Edge struct {
		a, b int
		d    int64
	}
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var a, b int
		var d int64
		fmt.Fscan(in, &a, &b, &d)
		edges[i] = Edge{a - 1, b - 1, d}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].d < edges[j].d })

	const LOG = 31
	adjMat := NewMatrix(n)
	pow := buildPowers(adjMat, LOG)
	curVec := NewBitset(n)
	curVec.Set(0)
	curTime := int64(0)
	target := n - 1
	idx := 0
	for idx < m {
		nextD := edges[idx].d
		diff := nextD - curTime
		if diff > 0 {
			if step, ok := earliestReach(curVec, pow, diff, target); ok {
				fmt.Println(curTime + step)
				return
			}
			curVec = applyPower(curVec, pow, diff)
			curTime = nextD
		}
		for idx < m && edges[idx].d == nextD {
			adjMat.AddEdge(edges[idx].a, edges[idx].b)
			idx++
		}
		pow = buildPowers(adjMat, LOG)
	}

	// Final search after all edges are added
	if curVec.Has(target) {
		fmt.Println(curTime)
		return
	}
	// Build adjacency list for BFS
	adjList := make([][]int, n)
	for i := 0; i < n; i++ {
		row := adjMat[i]
		for j := 0; j < n; j++ {
			if row.Has(j) {
				adjList[i] = append(adjList[i], j)
			}
		}
	}
	dist := bfsFromSet(curVec, adjList, target)
	if dist == -1 {
		fmt.Println("Impossible")
	} else {
		fmt.Println(curTime + int64(dist))
	}
}

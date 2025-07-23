package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxA = 200000

var spf [MaxA + 1]int

func sieve() {
	for i := 2; i <= MaxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= MaxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func divisors(x int) []int {
	res := []int{1}
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		m := len(res)
		pow := 1
		for c := 1; c <= cnt; c++ {
			pow *= p
			for i := 0; i < m; i++ {
				res = append(res, res[i]*pow)
			}
		}
	}
	return res
}

// DSU structure
type DSU struct {
	parent []int
	size   []int
}

func newDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) union(x, y int) {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	d.parent[ry] = rx
	d.size[rx] += d.size[ry]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	sieve()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	val := make([]int, n+1)
	maxV := 0
	nodesVal := make([][]int, MaxA+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &val[i])
		nodesVal[val[i]] = append(nodesVal[val[i]], i)
		if val[i] > maxV {
			maxV = val[i]
		}
	}

	type edge struct{ u, v int }
	edgesByG := make([][]edge, MaxA+1)
	edges := make([]edge, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v)
		g := gcd(val[edges[i].u], val[edges[i].v])
		edgesByG[g] = append(edgesByG[g], edges[i])
		if g > maxV {
			maxV = g
		}
	}

	// gather edges by multiples: for processing each d, we will iterate multiples m of d
	// Precompute prefix of nodes/edges for multiples using arrays of multiples indexes to avoid repeated loops? We'll not.

	cntDiv := make([]int64, maxV+1)

	for d := 1; d <= maxV; d++ {
		// gather nodes divisible by d
		nodes := make([]int, 0)
		for m := d; m <= maxV; m += d {
			nodes = append(nodes, nodesVal[m]...)
		}
		if len(nodes) == 0 {
			continue
		}
		id := make(map[int]int, len(nodes))
		for idx, node := range nodes {
			id[node] = idx
		}
		dsu := newDSU(len(nodes))
		// union edges
		for m := d; m <= maxV; m += d {
			for _, e := range edgesByG[m] {
				iu, ok1 := id[e.u]
				iv, ok2 := id[e.v]
				if ok1 && ok2 {
					dsu.union(iu, iv)
				}
			}
		}
		var count int64 = int64(len(nodes))
		for i := range dsu.parent {
			if dsu.parent[i] == i {
				s := int64(dsu.size[i])
				count += s * (s - 1) / 2
			}
		}
		cntDiv[d] = count
	}

	// compute exact gcd counts using inclusion-exclusion
	ans := make([]int64, maxV+1)
	for d := maxV; d >= 1; d-- {
		res := cntDiv[d]
		for m := d * 2; m <= maxV; m += d {
			res -= ans[m]
		}
		ans[d] = res
	}

	for i := 1; i <= maxV; i++ {
		if ans[i] > 0 {
			fmt.Fprintln(writer, i, ans[i])
		}
	}
}

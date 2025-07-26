package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

// DSU structure with range intersection
type DSU struct {
	parent []int
	size   []int
	L      []int64
	R      []int64
	len    []int64
}

func newDSU(n int, L, R []int64) *DSU {
	d := &DSU{
		parent: make([]int, n+1),
		size:   make([]int, n+1),
		L:      make([]int64, n+1),
		R:      make([]int64, n+1),
		len:    make([]int64, n+1),
	}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		d.L[i] = L[i]
		d.R[i] = R[i]
		if L[i] <= R[i] {
			d.len[i] = R[i] - L[i] + 1
		} else {
			d.len[i] = 0
		}
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

var prod int64
var valid bool = true

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	if valid {
		prod = prod * modInv(d.len[ra]%MOD) % MOD
		prod = prod * modInv(d.len[rb]%MOD) % MOD
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	if d.L[ra] < d.L[rb] {
		d.L[ra] = d.L[rb]
	}
	if d.R[ra] > d.R[rb] {
		d.R[ra] = d.R[rb]
	}
	if d.L[ra] <= d.R[ra] {
		d.len[ra] = d.R[ra] - d.L[ra] + 1
	} else {
		d.len[ra] = 0
	}
	if d.len[ra] == 0 {
		valid = false
		prod = 0
	} else if valid {
		prod = prod * (d.len[ra] % MOD) % MOD
	}
}

func getPath(parent []int, depth []int, u, v int) []int {
	var p1, p2 []int
	for depth[u] > depth[v] {
		p1 = append(p1, u)
		u = parent[u]
	}
	for depth[v] > depth[u] {
		p2 = append(p2, v)
		v = parent[v]
	}
	for u != v {
		p1 = append(p1, u)
		u = parent[u]
		p2 = append(p2, v)
		v = parent[v]
	}
	p1 = append(p1, u)
	for i := len(p2) - 1; i >= 0; i-- {
		p1 = append(p1, p2[i])
	}
	return p1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	parent := make([]int, n+1)
	parent[1] = 1
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &parent[i])
	}

	L := make([]int64, n+1)
	R := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &L[i], &R[i])
	}

	dsu := newDSU(n, L, R)
	prod = 1
	for i := 1; i <= n; i++ {
		if dsu.len[i] == 0 {
			valid = false
			prod = 0
			break
		}
		prod = prod * (dsu.len[i] % MOD) % MOD
	}

	depth := make([]int, n+1)
	// compute depth via parent chain (tree rooted at 1)
	for i := 2; i <= n; i++ {
		depth[i] = depth[parent[i]] + 1
	}

	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		if !valid {
			var a, b, c, d int
			fmt.Fscan(reader, &a, &b, &c, &d)
			fmt.Fprintln(writer, 0)
			continue
		}
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b, &c, &d)
		p1 := getPath(parent, depth, a, b)
		p2 := getPath(parent, depth, c, d)
		if len(p1) != len(p2) {
			// should not happen by problem statement
		}
		for j := 0; j < len(p1); j++ {
			dsu.union(p1[j], p2[j])
			if !valid {
				break
			}
		}
		if valid {
			fmt.Fprintln(writer, prod%MOD)
		} else {
			fmt.Fprintln(writer, 0)
		}
	}
}

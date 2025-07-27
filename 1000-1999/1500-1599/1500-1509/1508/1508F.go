package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU structure for range unions
type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	// permutation a is not needed for computing the answer
	for i := 0; i < n; i++ {
		var tmp int
		fmt.Fscan(reader, &tmp)
	}

	L := make([]int, q)
	R := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &L[i], &R[i])
		L[i]--
		R[i]--
	}

	dsu := NewDSU(n)
	// next array helps to merge ranges in amortized O(n)
	next := make([]int, n)
	for i := 0; i < n; i++ {
		next[i] = i
	}
	var getNext func(int) int
	getNext = func(x int) int {
		if next[x] != x {
			next[x] = getNext(next[x])
		}
		return next[x]
	}

	components := n
	for i := 0; i < q; i++ {
		l, r := L[i], R[i]
		cur := getNext(l)
		for cur < r {
			if dsu.union(cur, cur+1) {
				components--
			}
			next[cur] = getNext(cur + 1)
			cur = next[cur]
		}
		fmt.Fprintln(writer, n-components)
	}
}

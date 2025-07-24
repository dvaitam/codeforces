package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
	sets   []map[int]struct{}
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n+1),
		size:   make([]int, n+1),
		sets:   make([]map[int]struct{}, n+1),
	}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		m := make(map[int]struct{})
		m[i] = struct{}{}
		d.sets[i] = m
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b, time int, ans []int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	setA := d.sets[ra]
	setB := d.sets[rb]
	for v := range setB {
		if v < len(ans) {
			if _, ok := setA[v+1]; ok && ans[v] == 0 {
				ans[v] = time
			}
		}
		if v > 1 {
			if _, ok := setA[v-1]; ok && ans[v-1] == 0 {
				ans[v-1] = time
			}
		}
	}
	for v := range setB {
		setA[v] = struct{}{}
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.sets[rb] = nil
}

func buildSeg(arr []int) ([]int, int) {
	n := len(arr) - 1 // arr is 1-indexed
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int, size*2)
	for i := 1; i <= n; i++ {
		tree[size+i-1] = arr[i]
	}
	for i := size - 1; i >= 1; i-- {
		if tree[i<<1] > tree[i<<1|1] {
			tree[i] = tree[i<<1]
		} else {
			tree[i] = tree[i<<1|1]
		}
	}
	return tree, size
}

func querySeg(tree []int, size int, l, r int) int {
	if l > r {
		return 0
	}
	l += size - 1
	r += size - 1
	res := 0
	for l <= r {
		if l&1 == 1 {
			if tree[l] > res {
				res = tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if tree[r] > res {
				res = tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, q int
		fmt.Fscan(reader, &n, &m, &q)
		edgesU := make([]int, m)
		edgesV := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &edgesU[i], &edgesV[i])
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &queries[i][0], &queries[i][1])
		}

		ansPair := make([]int, n)
		dsu := NewDSU(n)
		for i := 0; i < m; i++ {
			dsu.union(edgesU[i], edgesV[i], i+1, ansPair)
		}

		tree, size := buildSeg(ansPair)
		for i := 0; i < q; i++ {
			l := queries[i][0]
			r := queries[i][1]
			if l == r {
				fmt.Fprint(writer, 0)
			} else {
				res := querySeg(tree, size, l, r-1)
				fmt.Fprint(writer, res)
			}
			if i+1 < q {
				fmt.Fprint(writer, " ")
			}
		}
		if t > 1 {
			fmt.Fprintln(writer)
		}
	}
}

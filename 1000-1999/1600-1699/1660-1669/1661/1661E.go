package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	comp   int
	parent [6]int16
	active [6]bool
}

func find(parent []int16, x int) int {
	for int(parent[x]) != x {
		x = int(parent[x])
	}
	return x
}

func union(parent []int16, x, y int) bool {
	fx := find(parent, x)
	fy := find(parent, y)
	if fx == fy {
		return false
	}
	parent[fy] = int16(fx)
	return true
}

func buildLeaf(cols []string, idx int) Node {
	var n Node
	for i := 0; i < 6; i++ {
		n.parent[i] = int16(i)
	}
	// process each row
	for r := 0; r < 3; r++ {
		if cols[r][idx] == '1' {
			n.active[r] = true
			n.active[r+3] = true
			union(n.parent[:], r, r+3)
			if r > 0 && cols[r-1][idx] == '1' {
				union(n.parent[:], r-1, r)
				union(n.parent[:], r-1+3, r+3)
			}
		}
	}
	vis := make(map[int]bool)
	for r := 0; r < 3; r++ {
		if n.active[r] {
			rt := find(n.parent[:], r)
			if !vis[rt] {
				vis[rt] = true
				n.comp++
			}
		}
	}
	return n
}

func merge(a, b Node) Node {
	parent := make([]int16, 12)
	for i := range parent {
		parent[i] = int16(i)
	}
	active := make([]bool, 12)
	for i := 0; i < 6; i++ {
		active[i] = a.active[i]
		active[6+i] = b.active[i]
	}
	// replicate unions within a
	for i := 0; i < 6; i++ {
		if !active[i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !active[j] {
				continue
			}
			if find(a.parent[:], i) == find(a.parent[:], j) {
				union(parent, i, j)
			}
		}
	}
	// replicate unions within b
	for i := 0; i < 6; i++ {
		if !active[6+i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !active[6+j] {
				continue
			}
			if find(b.parent[:], i) == find(b.parent[:], j) {
				union(parent, 6+i, 6+j)
			}
		}
	}
	comp := a.comp + b.comp
	for r := 0; r < 3; r++ {
		if a.active[3+r] && b.active[r] {
			if union(parent, 3+r, 6+r) {
				comp--
			}
		}
	}
	var res Node
	res.comp = comp
	for i := 0; i < 6; i++ {
		res.parent[i] = int16(i)
	}
	for i := 0; i < 3; i++ {
		res.active[i] = a.active[i]
	}
	for i := 0; i < 3; i++ {
		res.active[3+i] = b.active[3+i]
	}
	idxMap := []int{0, 1, 2, 9, 10, 11}
	for i := 0; i < 6; i++ {
		if !res.active[i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !res.active[j] {
				continue
			}
			if find(parent, idxMap[i]) == find(parent, idxMap[j]) {
				union(res.parent[:], i, j)
			}
		}
	}
	return res
}

var seg []Node
var cols []string

func build(id, l, r int) {
	if l == r {
		seg[id] = buildLeaf(cols, l)
		return
	}
	mid := (l + r) / 2
	build(id*2, l, mid)
	build(id*2+1, mid+1, r)
	seg[id] = merge(seg[id*2], seg[id*2+1])
}

func query(id, l, r, L, R int) Node {
	if L <= l && r <= R {
		return seg[id]
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(id*2, l, mid, L, R)
	}
	if L > mid {
		return query(id*2+1, mid+1, r, L, R)
	}
	left := query(id*2, l, mid, L, mid)
	right := query(id*2+1, mid+1, r, mid+1, R)
	return merge(left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cols = make([]string, 3)
	for i := 0; i < 3; i++ {
		fmt.Fscan(in, &cols[i])
	}
	seg = make([]Node, 4*n)
	build(1, 0, n-1)
	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		r--
		res := query(1, 0, n-1, l, r)
		fmt.Fprintln(out, res.comp)
	}
}

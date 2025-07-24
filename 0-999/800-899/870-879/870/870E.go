package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type pair struct{ x, y int }
	pts := make([]pair, n)
	xMap := make(map[int]int)
	yMap := make(map[int]int)
	xi := 0
	yi := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
		if _, ok := xMap[pts[i].x]; !ok {
			xMap[pts[i].x] = xi
			xi++
		}
		if _, ok := yMap[pts[i].y]; !ok {
			yMap[pts[i].y] = yi
			yi++
		}
	}

	total := xi + yi
	dsu := NewDSU(total)
	edges := make([][2]int, n)
	for i := 0; i < n; i++ {
		xID := xMap[pts[i].x]
		yID := yMap[pts[i].y] + xi
		dsu.Union(xID, yID)
		edges[i] = [2]int{xID, yID}
	}

	edgeCnt := make([]int, total)
	for _, e := range edges {
		root := dsu.Find(e[0])
		edgeCnt[root]++
	}

	vertCnt := make([]int, total)
	for i := 0; i < total; i++ {
		root := dsu.Find(i)
		vertCnt[root]++
	}

	pow2 := make([]int64, total+1)
	pow2[0] = 1
	for i := 1; i <= total; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}

	seen := make(map[int]bool)
	ans := int64(1)
	for i := 0; i < total; i++ {
		root := dsu.Find(i)
		if seen[root] {
			continue
		}
		seen[root] = true
		v := vertCnt[root]
		e := edgeCnt[root]
		if e == v-1 {
			ans = ans * ((pow2[v] - 1 + mod) % mod) % mod
		} else {
			ans = ans * pow2[v] % mod
		}
	}

	fmt.Fprintln(out, ans)
}

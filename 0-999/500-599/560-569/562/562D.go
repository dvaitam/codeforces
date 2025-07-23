package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU data structure with path compression and union by rank
type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	return &DSU{parent: p, rank: r}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
}

// helper structure to skip already merged positions during range unions
var nextPos []int

func getNext(x int) int {
	if nextPos[x] != x {
		nextPos[x] = getNext(nextPos[x])
	}
	return nextPos[x]
}

func unionRange(dsu *DSU, l, r int) {
	cur := getNext(l)
	for cur < r {
		dsu.Union(cur, cur+1)
		nextPos[cur] = getNext(cur + 1)
		cur = getNext(cur)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	dsu := NewDSU(n)
	nextPos = make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		nextPos[i] = i
	}

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(in, &t, &x, &y)
		if t == 1 {
			dsu.Union(x, y)
		} else if t == 2 {
			if x > y {
				x, y = y, x
			}
			unionRange(dsu, x, y)
		} else if t == 3 {
			if dsu.Find(x) == dsu.Find(y) {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

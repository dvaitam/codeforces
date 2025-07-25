package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return
	}
	if d.rank[x] < d.rank[y] {
		x, y = y, x
	}
	d.parent[y] = x
	if d.rank[x] == d.rank[y] {
		d.rank[x]++
	}
}

func canSort(p []int, l, r int) bool {
	n := len(p)
	d := NewDSU(n + 1)
	for x := 1; x <= n; x++ {
		for y := x + 1; y <= n; y++ {
			s := x + y
			if s >= l && s <= r {
				d.Union(x, y)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if d.Find(i) != d.Find(p[i-1]) {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		ans := 0
		maxSum := 2 * n
		for l := 1; l <= maxSum; l++ {
			for r := l; r <= maxSum; r++ {
				if canSort(p, l, r) {
					ans++
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

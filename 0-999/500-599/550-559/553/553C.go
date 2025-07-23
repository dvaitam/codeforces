package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1000000007

type DSU struct {
	parent []int
	parity []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	parity := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		parity[i] = 0
	}
	return &DSU{parent: parent, parity: parity}
}

func (d *DSU) Find(x int) (int, int) {
	if d.parent[x] != x {
		root, par := d.Find(d.parent[x])
		d.parent[x] = root
		d.parity[x] ^= par
	}
	return d.parent[x], d.parity[x]
}

func (d *DSU) Union(x, y, val int) bool {
	rx, px := d.Find(x)
	ry, py := d.Find(y)
	if rx == ry {
		return (px ^ py) == val
	}
	d.parent[rx] = ry
	d.parity[rx] = px ^ py ^ val
	return true
}

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	dsu := NewDSU(n)
	ok := true
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if !dsu.Union(a, b, c^1) {
			ok = false
		}
	}
	if !ok {
		fmt.Println(0)
		return
	}
	seen := make(map[int]bool)
	comp := 0
	for i := 1; i <= n; i++ {
		r, _ := dsu.Find(i)
		if !seen[r] {
			seen[r] = true
			comp++
		}
	}
	ans := modPow(2, comp-1)
	fmt.Println(ans)
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	diff   []int64
	mod    int64
}

func NewDSU(n int, mod int64) *DSU {
	parent := make([]int, n)
	diff := make([]int64, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		diff[i] = 0
	}
	return &DSU{parent: parent, diff: diff, mod: mod}
}

func (d *DSU) find(x int) (int, int64) {
	if d.parent[x] == x {
		return x, 0
	}
	r, w := d.find(d.parent[x])
	d.diff[x] = (d.diff[x] + w) % d.mod
	d.parent[x] = r
	return r, d.diff[x]
}

func (d *DSU) unite(x, y int, val int64) bool { // enforce w[x]-w[y]=val (mod mod)
	rx, dx := d.find(x)
	ry, dy := d.find(y)
	if rx == ry {
		if (dx-dy-val)%d.mod != 0 {
			return false
		}
		return true
	}
	d.parent[rx] = ry
	d.diff[rx] = ((val-dx+dy)%d.mod + d.mod) % d.mod
	return true
}

const MOD int64 = 1_000_000_007

func powMod(a, b, mod int64) int64 {
	res := int64(1)
	a %= mod
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var tc int
	fmt.Fscan(in, &tc)
	for ; tc > 0; tc-- {
		var n, m int
		var h int64
		fmt.Fscan(in, &n, &m, &h)
		dsu := NewDSU(n+m, h)
		ok := true
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int64
				fmt.Fscan(in, &x)
				if !ok {
					continue
				}
				if x != -1 {
					if !dsu.unite(i, n+j, x%h) {
						ok = false
					}
				}
			}
		}
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}
		comp := 0
		for i := 0; i < n+m; i++ {
			r, _ := dsu.find(i)
			if r == i {
				comp++
			}
		}
		ans := powMod(h, int64(comp-1), MOD)
		fmt.Fprintln(out, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000009

type DSU struct {
	p []int
	s []int
}

func newDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1), s: make([]int, n+1)}
	for i := 0; i <= n; i++ {
		d.p[i] = i
		d.s[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.s[ra] < d.s[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.s[ra] += d.s[rb]
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	dsu := newDSU(n)
	pow := int64(1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if !dsu.union(u, v) {
			pow = (pow * 2) % MOD
		}
		ans := (pow - 1 + MOD) % MOD
		fmt.Fprintln(out, ans)
	}
}

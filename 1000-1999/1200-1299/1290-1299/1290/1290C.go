package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU structure with parity and cost tracking
type DSU struct {
	parent []int
	parity []int
	size   []int
	cnt    [][2]int // cnt[x][0]: color 0 count, cnt[x][1]: color 1 count
	forced []int    // -1 if no constraint, otherwise 0 or 1
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		parity: make([]int, n),
		size:   make([]int, n),
		cnt:    make([][2]int, n),
		forced: make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		d.forced[i] = -1
	}
	d.forced[0] = 0 // node 0 is constant 0
	for i := 1; i < n; i++ {
		d.cnt[i][0] = 1 // each set initially has color 0
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		px := d.parent[x]
		d.parent[x] = d.find(px)
		d.parity[x] ^= d.parity[px]
	}
	return d.parent[x]
}

func (d *DSU) cost(x int) int {
	if d.forced[x] == -1 {
		if d.cnt[x][0] < d.cnt[x][1] {
			return d.cnt[x][0]
		}
		return d.cnt[x][1]
	}
	if d.forced[x] == 0 {
		return d.cnt[x][1]
	}
	return d.cnt[x][0]
}

func (d *DSU) unite(a, b, w int, ans *int) {
	ra := d.find(a)
	rb := d.find(b)
	w ^= d.parity[a] ^ d.parity[b]
	if ra == rb {
		return
	}
	if d.size[ra] > d.size[rb] {
		ra, rb = rb, ra
	}
	*ans -= d.cost(ra)
	*ans -= d.cost(rb)
	d.parent[ra] = rb
	d.parity[ra] = w
	d.size[rb] += d.size[ra]
	c0 := d.cnt[ra][0]
	c1 := d.cnt[ra][1]
	if w == 1 {
		c0, c1 = c1, c0
	}
	d.cnt[rb][0] += c0
	d.cnt[rb][1] += c1
	if d.forced[ra] != -1 {
		val := d.forced[ra] ^ w
		if d.forced[rb] == -1 {
			d.forced[rb] = val
		}
	}
	*ans += d.cost(rb)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	pos := make([][]int, n)
	for i := 1; i <= k; i++ {
		var c int
		fmt.Fscan(in, &c)
		for j := 0; j < c; j++ {
			var x int
			fmt.Fscan(in, &x)
			pos[x-1] = append(pos[x-1], i)
		}
	}

	dsu := NewDSU(k + 1)
	ans := 0
	for i := 0; i < n; i++ {
		v := int('1' - s[i]) // 1 - s[i]
		if len(pos[i]) == 1 {
			dsu.unite(pos[i][0], 0, v, &ans)
		} else if len(pos[i]) == 2 {
			dsu.unite(pos[i][0], pos[i][1], v, &ans)
		}
		fmt.Fprintln(out, ans)
	}
}

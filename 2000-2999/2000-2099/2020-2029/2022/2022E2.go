package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod      = 1000000007
	maxNodes = 200000 + 5
	maxExp   = 30 * maxNodes
)

var pow2 []int

type dsu struct {
	parent      []int
	xorToParent []int
	size        []int
	comps       int
}

func newDSU(n int) *dsu {
	parent := make([]int, n)
	xorToParent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &dsu{parent: parent, xorToParent: xorToParent, size: size, comps: n}
}

func (d *dsu) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, 0
	}
	root, xr := d.find(d.parent[x])
	total := xr ^ d.xorToParent[x]
	d.parent[x] = root
	d.xorToParent[x] = total
	return root, total
}

func (d *dsu) union(a, b, val int) bool {
	ra, xa := d.find(a)
	rb, xb := d.find(b)
	if ra == rb {
		return (xa ^ xb) == val
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
		xa, xb = xb, xa
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.xorToParent[rb] = xa ^ xb ^ val
	d.comps--
	return true
}

func ensurePow2() {
	if pow2 != nil {
		return
	}
	pow2 = make([]int, maxExp+1)
	pow2[0] = 1
	for i := 1; i <= maxExp; i++ {
		pow2[i] = (pow2[i-1] << 1) % mod
	}
}

func main() {
	ensurePow2()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k, q int
		fmt.Fscan(in, &n, &m, &k, &q)
		totalNodes := n + m
		d := newDSU(totalNodes)
		inconsistent := false
		readConstraint := func() (int, int, int) {
			var r, c, v int
			fmt.Fscan(in, &r, &c, &v)
			return r - 1, c - 1, v
		}
		for i := 0; i < k; i++ {
			r, c, v := readConstraint()
			row := r
			col := n + c
			if !inconsistent {
				if !d.union(row, col, v) {
					inconsistent = true
				}
			}
		}
		writeAnswer := func() {
			if inconsistent {
				fmt.Fprintln(out, 0)
				return
			}
			exp := 30 * (d.comps - 1)
			fmt.Fprintln(out, pow2[exp])
		}
		writeAnswer()
		for i := 0; i < q; i++ {
			r, c, v := readConstraint()
			row := r
			col := n + c
			if !inconsistent {
				if !d.union(row, col, v) {
					inconsistent = true
				}
			}
			writeAnswer()
		}
	}
}

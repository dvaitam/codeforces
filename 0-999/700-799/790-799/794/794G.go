package main

import (
	"bufio"
	"fmt"
	"os"
)

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x != y {
		d.parent[x] = y
	}
}

const mod = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

// compute number of solutions for fixed strings x, y consisting of 'A' and 'B'
// and given lengths a = |s|, b = |t|.
func check(x, y string, a, b int) (bool, int) {
	d := newDSU(a + b)
	i, j := 0, 0
	offX, offY := 0, 0
	for i < len(x) && j < len(y) {
		chX := x[i]
		chY := y[j]
		lenX := a
		if chX == 'B' {
			lenX = b
		}
		lenY := a
		if chY == 'B' {
			lenY = b
		}
		step := lenX - offX
		if lenY-offY < step {
			step = lenY - offY
		}
		for k := 0; k < step; k++ {
			idx1 := offX + k
			if chX == 'B' {
				idx1 = a + offX + k
			}
			idx2 := offY + k
			if chY == 'B' {
				idx2 = a + offY + k
			}
			d.union(idx1, idx2)
		}
		offX += step
		offY += step
		if offX == lenX {
			offX = 0
			i++
		}
		if offY == lenY {
			offY = 0
			j++
		}
	}
	if i != len(x) || j != len(y) || offX != 0 || offY != 0 {
		return false, 0
	}
	comp := 0
	for idx := 0; idx < a+b; idx++ {
		if d.find(idx) == idx {
			comp++
		}
	}
	return true, comp
}

func solve(x, y string, n int) int64 {
	ans := int64(0)
	for a := 1; a <= n; a++ {
		for b := 1; b <= n; b++ {
			ok, comp := check(x, y, a, b)
			if ok {
				ans = (ans + modPow(2, int64(comp))) % mod
			}
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var c, d string
	var n int
	if _, err := fmt.Fscan(in, &c); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &d); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// This solution does not handle '?' characters efficiently.
	// It simply assumes the strings contain only 'A' and 'B'.
	// The algorithm is intended for small input sizes.
	fmt.Println(solve(c, d, n))
}

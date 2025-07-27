package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n: n, bit: make([]int64, n+2)} }
func (f *Fenwick) Add(idx int, delta int64) {
	idx++
	for idx < len(f.bit) {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}
func (f *Fenwick) Sum(idx int) int64 {
	idx++
	res := int64(0)
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}
func (f *Fenwick) RangeSum(l, r int) int64 {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	const MAX = 300000
	cnt := NewFenwick(MAX + 2)
	sum := NewFenwick(MAX + 2)
	divf := NewFenwick(MAX + 2)
	ans := make([]int64, n)
	var acc int64
	for i, x := range a {
		var part1 int64
		for m := 0; m <= MAX; m += x {
			r := m + x - 1
			if r > MAX {
				r = MAX
			}
			c := cnt.RangeSum(m, r)
			s := sum.RangeSum(m, r)
			part1 += s - c*int64(m)
		}
		tot := int64(i)
		part2 := int64(x)*tot - divf.Sum(x)
		acc += part1 + part2
		ans[i] = acc
		cnt.Add(x, 1)
		sum.Add(x, int64(x))
		for m := x; m <= MAX; m += x {
			divf.Add(m, int64(x))
		}
	}
	w := bufio.NewWriter(os.Stdout)
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, v)
	}
	fmt.Fprintln(w)
	w.Flush()
}

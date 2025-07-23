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

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) Add(i int, delta int64) {
	for i <= f.n {
		f.bit[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int64 {
	if i > f.n {
		i = f.n
	}
	s := int64(0)
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, a, b, q int
	if _, err := fmt.Fscan(in, &n, &k, &a, &b, &q); err != nil {
		return
	}
	// orders count per day
	cnt := make([]int, n+1)
	// Fenwick trees for truncated counts
	ftA := NewFenwick(n)
	ftB := NewFenwick(n)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var d, add int
			fmt.Fscan(in, &d, &add)
			old := cnt[d]
			cnt[d] += add
			newVal := cnt[d]
			deltaA := int64(min(newVal, a) - min(old, a))
			deltaB := int64(min(newVal, b) - min(old, b))
			if deltaA != 0 {
				ftA.Add(d, deltaA)
			}
			if deltaB != 0 {
				ftB.Add(d, deltaB)
			}
		} else if t == 2 {
			var p int
			fmt.Fscan(in, &p)
			before := ftB.Sum(p - 1)
			after := ftA.Sum(n) - ftA.Sum(p+k-1)
			fmt.Fprintln(out, before+after)
		}
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

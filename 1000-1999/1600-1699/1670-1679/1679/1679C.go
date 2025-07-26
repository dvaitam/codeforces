package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, delta int) {
	for i <= f.n {
		f.bit[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	if i > f.n {
		i = f.n
	}
	s := 0
	for i > 0 {
		s += f.bit[i]
		i &= i - 1
	}
	return s
}

func (f *Fenwick) RangeSum(l, r int) int {
	if l > r {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	rowCount := make([]int, n+1)
	colCount := make([]int, n+1)
	rowBit := NewFenwick(n)
	colBit := NewFenwick(n)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x, y int
			fmt.Fscan(in, &x, &y)
			rowCount[x]++
			if rowCount[x] == 1 {
				rowBit.Add(x, 1)
			}
			colCount[y]++
			if colCount[y] == 1 {
				colBit.Add(y, 1)
			}
		} else if t == 2 {
			var x, y int
			fmt.Fscan(in, &x, &y)
			rowCount[x]--
			if rowCount[x] == 0 {
				rowBit.Add(x, -1)
			}
			colCount[y]--
			if colCount[y] == 0 {
				colBit.Add(y, -1)
			}
		} else if t == 3 {
			var x1, y1, x2, y2 int
			fmt.Fscan(in, &x1, &y1, &x2, &y2)
			rows := rowBit.RangeSum(x1, x2)
			cols := colBit.RangeSum(y1, y2)
			if rows == x2-x1+1 || cols == y2-y1+1 {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		}
	}
}

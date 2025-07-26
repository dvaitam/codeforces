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

func (f *Fenwick) add(idx int, val int64) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) sum(idx int) int64 {
	var res int64
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

type RangeFenwick struct {
	n      int
	b1, b2 *Fenwick
}

func NewRangeFenwick(n int) *RangeFenwick {
	return &RangeFenwick{n: n, b1: NewFenwick(n + 2), b2: NewFenwick(n + 2)}
}

func (rf *RangeFenwick) rangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	rf.b1.add(l, val)
	rf.b1.add(r+1, -val)
	rf.b2.add(l, val*int64(l-1))
	rf.b2.add(r+1, -val*int64(r))
}

func (rf *RangeFenwick) prefixSum(x int) int64 {
	return rf.b1.sum(x)*int64(x) - rf.b2.sum(x)
}

type event struct {
	l, r int
}

type query struct {
	r   int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	p := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
		pos[p[i]] = i
	}

	left := make([]int, n+1)
	stack := []int{}
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			left[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	right := make([]int, n+1)
	stack = []int{}
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			right[i] = stack[len(stack)-1]
		} else {
			right[i] = n + 1
		}
		stack = append(stack, i)
	}

	add := make([][]event, n+2)
	rem := make([][]event, n+2)

	for a := 1; a <= n; a++ {
		for b := a + 1; a*b <= n; b++ {
			m := a * b
			ia, ib, im := pos[a], pos[b], pos[m]
			lb := left[im] + 1
			rb := right[im] - 1
			minPos := ia
			if ib < minPos {
				minPos = ib
			}
			if im < minPos {
				minPos = im
			}
			maxPos := ia
			if ib > maxPos {
				maxPos = ib
			}
			if im > maxPos {
				maxPos = im
			}
			if lb <= minPos && maxPos <= rb {
				l1 := lb
				l2 := minPos
				r1 := maxPos
				r2 := rb
				if l1 <= l2 && r1 <= r2 {
					add[l2] = append(add[l2], event{r1, r2})
					if l1 > 1 {
						rem[l1-1] = append(rem[l1-1], event{r1, r2})
					}
				}
			}
		}
	}

	queries := make([][]query, n+1)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		queries[l] = append(queries[l], query{r, i})
	}

	rf := NewRangeFenwick(n + 2)
	ans := make([]int64, q)

	for l := n; l >= 1; l-- {
		for _, e := range add[l] {
			rf.rangeAdd(e.l, e.r, 1)
		}
		for _, e := range rem[l] {
			rf.rangeAdd(e.l, e.r, -1)
		}
		for _, qu := range queries[l] {
			ans[qu.idx] = rf.prefixSum(qu.r)
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}

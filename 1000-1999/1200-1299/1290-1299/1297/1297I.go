package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxD = 100000

type node struct {
	stack []int
	mx    int
}

var tree []node
var L, R []int

func pushStack(idx, val int) {
	tree[idx].stack = append(tree[idx].stack, val)
	if val > tree[idx].mx {
		tree[idx].mx = val
	}
}

func pull(idx int) {
	mx := 0
	st := tree[idx].stack
	if len(st) > 0 {
		mx = st[len(st)-1]
	}
	l := idx << 1
	if l < len(tree) {
		if tree[l].mx > mx {
			mx = tree[l].mx
		}
		if tree[l+1].mx > mx {
			mx = tree[l+1].mx
		}
	}
	tree[idx].mx = mx
}

func pushDown(idx int) {
	if len(tree[idx].stack) == 0 {
		return
	}
	l := idx << 1
	if l >= len(tree) {
		return
	}
	for _, v := range tree[idx].stack {
		pushStack(l, v)
		pushStack(l+1, v)
	}
	tree[idx].stack = tree[idx].stack[:0]
	pull(idx)
}

func updateAdd(idx, l, r, Lq, Rq, val int) {
	if Lq <= l && r <= Rq {
		pushStack(idx, val)
		return
	}
	pushDown(idx)
	mid := (l + r) >> 1
	if Lq <= mid {
		updateAdd(idx<<1, l, mid, Lq, Rq, val)
	}
	if Rq > mid {
		updateAdd(idx<<1|1, mid+1, r, Lq, Rq, val)
	}
	pull(idx)
}

func updateRemove(idx, l, r, Lq, Rq int) {
	if Lq <= l && r <= Rq {
		st := tree[idx].stack
		tree[idx].stack = st[:len(st)-1]
		pull(idx)
		return
	}
	pushDown(idx)
	mid := (l + r) >> 1
	if Lq <= mid {
		updateRemove(idx<<1, l, mid, Lq, Rq)
	}
	if Rq > mid {
		updateRemove(idx<<1|1, mid+1, r, Lq, Rq)
	}
	pull(idx)
}

func query(idx, l, r, Lq, Rq int) int {
	if Lq <= l && r <= Rq {
		return tree[idx].mx
	}
	pushDown(idx)
	mid := (l + r) >> 1
	res := 0
	if Lq <= mid {
		if v := query(idx<<1, l, mid, Lq, Rq); v > res {
			res = v
		}
	}
	if Rq > mid {
		if v := query(idx<<1|1, mid+1, r, Lq, Rq); v > res {
			res = v
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, d int
	if _, err := fmt.Fscan(reader, &n, &d); err != nil {
		return
	}

	tree = make([]node, 4*d+5)
	L = make([]int, n+1)
	R = make([]int, n+1)

	alive := 0
	for i := 1; i <= n; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		L[i] = l
		R[i] = r
		// remove covered blocks
		for {
			j := query(1, 1, d, l, r)
			if j == 0 {
				break
			}
			lj, rj := L[j], R[j]
			if l <= lj && rj <= r {
				updateRemove(1, 1, d, lj, rj)
				alive--
			} else {
				break
			}
		}
		updateAdd(1, 1, d, l, r, i)
		alive++
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, alive)
	}
	fmt.Fprintln(writer)
}

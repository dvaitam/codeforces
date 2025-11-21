package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = -1 << 60

type segTree struct {
	n      int
	maxVal []int
	maxIdx []int
	lazy   []int
}

func newSegTree(arr []int) *segTree {
	n := len(arr)
	st := &segTree{
		n:      n,
		maxVal: make([]int, 4*n),
		maxIdx: make([]int, 4*n),
		lazy:   make([]int, 4*n),
	}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segTree) build(node, l, r int, arr []int) {
	if l == r {
		st.maxVal[node] = arr[l]
		st.maxIdx[node] = l
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, arr)
	st.build(node<<1|1, mid+1, r, arr)
	st.pull(node)
}

func (st *segTree) pull(node int) {
	left, right := node<<1, node<<1|1
	if st.maxVal[left] > st.maxVal[right] {
		st.maxVal[node] = st.maxVal[left]
		st.maxIdx[node] = st.maxIdx[left]
	} else if st.maxVal[right] > st.maxVal[left] {
		st.maxVal[node] = st.maxVal[right]
		st.maxIdx[node] = st.maxIdx[right]
	} else {
		st.maxVal[node] = st.maxVal[left]
		if st.maxIdx[left] < st.maxIdx[right] {
			st.maxIdx[node] = st.maxIdx[left]
		} else {
			st.maxIdx[node] = st.maxIdx[right]
		}
	}
}

func (st *segTree) apply(node, val int) {
	st.maxVal[node] += val
	st.lazy[node] += val
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		st.apply(node<<1, st.lazy[node])
		st.apply(node<<1|1, st.lazy[node])
		st.lazy[node] = 0
	}
}

func (st *segTree) addRange(l, r, val int) {
	if l > r {
		return
	}
	st.add(1, 0, st.n-1, l, r, val)
}

func (st *segTree) add(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(node, val)
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.add(node<<1, l, mid, ql, qr, val)
	st.add(node<<1|1, mid+1, r, ql, qr, val)
	st.pull(node)
}

func (st *segTree) queryMax() (int, int) {
	return st.maxVal[1], st.maxIdx[1]
}

func (st *segTree) queryPoint(pos int) int {
	return st.get(1, 0, st.n-1, pos)
}

func (st *segTree) get(node, l, r, pos int) int {
	if l == r {
		return st.maxVal[node]
	}
	st.push(node)
	mid := (l + r) >> 1
	if pos <= mid {
		return st.get(node<<1, l, mid, pos)
	}
	return st.get(node<<1|1, mid+1, r, pos)
}

func (st *segTree) setNegInf(pos int) {
	cur := st.queryPoint(pos)
	st.addRange(pos, pos, int(negInf)-cur)
}

func solveCase(s string) string {
	n := len(s)
	digits := make([]int, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = int(s[i] - '0')
		scores[i] = digits[i] - i
	}
	st := newSegTree(scores)
	removed := make([]bool, n)
	res := make([]byte, 0, n)
	for {
		val, idx := st.queryMax()
		if val < 0 {
			break
		}
		res = append(res, byte('0'+val))
		removed[idx] = true
		st.setNegInf(idx)
		if idx+1 < n {
			st.addRange(idx+1, n-1, 1)
		}
	}
	for i := 0; i < n; i++ {
		if !removed[i] {
			res = append(res, byte('0'+digits[i]))
		}
	}
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solveCase(s))
	}
}

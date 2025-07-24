package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type SegTree struct {
	n    int
	tree []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1 // arr is 1-based
	st := &SegTree{n: n, tree: make([]int, 4*n)}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int) {
	if l == r {
		st.tree[node] = arr[l]
		return
	}
	m := (l + r) / 2
	st.build(node*2, l, m, arr)
	st.build(node*2+1, m+1, r, arr)
	st.tree[node] = st.tree[node*2] | st.tree[node*2+1]
}

func (st *SegTree) Update(node, l, r, pos, val int) {
	if l == r {
		st.tree[node] = val
		return
	}
	m := (l + r) / 2
	if pos <= m {
		st.Update(node*2, l, m, pos, val)
	} else {
		st.Update(node*2+1, m+1, r, pos, val)
	}
	st.tree[node] = st.tree[node*2] | st.tree[node*2+1]
}

func (st *SegTree) Query(node, l, r, L, R int) int {
	if L <= l && r <= R {
		return st.tree[node]
	}
	m := (l + r) / 2
	res := 0
	if L <= m {
		res |= st.Query(node*2, l, m, L, R)
	}
	if R > m {
		res |= st.Query(node*2+1, m+1, r, L, R)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = 1 << (s[i-1] - 'a')
	}
	st := NewSegTree(arr)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var pos int
			var c string
			fmt.Fscan(reader, &pos, &c)
			ch := c[0]
			if s[pos-1] != ch {
				sBytes := []byte(s)
				sBytes[pos-1] = ch
				s = string(sBytes)
				st.Update(1, 1, n, pos, 1<<(ch-'a'))
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			mask := st.Query(1, 1, n, l, r)
			fmt.Fprintln(writer, bits.OnesCount(uint(mask)))
		}
	}
}

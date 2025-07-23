package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9 + 7)

type SegTree struct {
	n    int
	tree []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	st := &SegTree{n: size, tree: make([]int, size*2)}
	for i := range st.tree {
		st.tree[i] = INF
	}
	return st
}

func (st *SegTree) Update(pos, val int) {
	pos += st.n
	st.tree[pos] = val
	for pos > 1 {
		pos >>= 1
		if st.tree[pos<<1] < st.tree[pos<<1|1] {
			st.tree[pos] = st.tree[pos<<1]
		} else {
			st.tree[pos] = st.tree[pos<<1|1]
		}
	}
}

// Query returns the smallest index i such that last[i] < val, or -1 if none
func (st *SegTree) Query(val int) int {
	if st.tree[1] >= val {
		return -1
	}
	idx, l, r := 1, 0, st.n-1
	for l != r {
		mid := (l + r) >> 1
		if st.tree[idx<<1] < val {
			idx = idx << 1
			r = mid
		} else {
			idx = idx<<1 | 1
			l = mid + 1
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	st := NewSegTree(n)
	sequences := make([][]int, 0)

	for _, x := range arr {
		idx := st.Query(x)
		if idx == -1 || idx >= len(sequences) {
			sequences = append(sequences, []int{x})
			st.Update(len(sequences)-1, x)
		} else {
			sequences[idx] = append(sequences[idx], x)
			st.Update(idx, x)
		}
	}

	for _, seq := range sequences {
		for j, v := range seq {
			if j > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

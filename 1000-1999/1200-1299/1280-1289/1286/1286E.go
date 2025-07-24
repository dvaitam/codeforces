package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// Segment tree for range minimum queries
// 1-indexed positions

type SegTree struct {
	n    int
	size int
	tree []int64
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int64, 2*size)
	for i := range tree {
		tree[i] = (1<<63 - 1) // INF
	}
	return &SegTree{n: n, size: size, tree: tree}
}

func (st *SegTree) Update(pos int, val int64) {
	i := pos + st.size - 1
	st.tree[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if st.tree[i<<1] < st.tree[i<<1|1] {
			st.tree[i] = st.tree[i<<1]
		} else {
			st.tree[i] = st.tree[i<<1|1]
		}
	}
}

func (st *SegTree) Query(l, r int) int64 {
	l += st.size - 1
	r += st.size - 1
	res := int64(1<<63 - 1)
	for l <= r {
		if l&1 == 1 {
			if st.tree[l] < res {
				res = st.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.tree[r] < res {
				res = st.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	S := make([]byte, 0, n)
	W := make([]int64, 0, n+1)
	pi := make([]int, n+1)
	seg := NewSegTree(n + 2)

	MASK := int64((1 << 30) - 1)
	ans := big.NewInt(0)
	prevAns := big.NewInt(0)

	for i := 1; i <= n; i++ {
		var ch string
		var w int64
		fmt.Fscan(in, &ch, &w)

		// decrypt using previous answer
		shift := new(big.Int).Mod(prevAns, big.NewInt(26))
		shiftInt := int(shift.Int64())
		c := ch[0]
		c = byte('a' + (int(c-'a')+shiftInt)%26)
		maskVal := new(big.Int).And(prevAns, big.NewInt(MASK))
		w ^= maskVal.Int64()

		S = append(S, c)
		W = append(W, w)
		seg.Update(i, w)

		// compute prefix function at position i (0-indexed on S)
		j := pi[i-1]
		for j > 0 && S[i-1] != S[j] {
			j = pi[j]
		}
		if S[i-1] == S[j] {
			j++
		}
		pi[i] = j

		// traverse borders and accumulate
		cur := i
		for cur > 0 {
			l := cur
			minv := seg.Query(i-l+1, i)
			ans.Add(ans, big.NewInt(minv))
			cur = pi[cur]
		}

		fmt.Fprintln(out, ans.String())
		prevAns.Set(ans)
	}
}

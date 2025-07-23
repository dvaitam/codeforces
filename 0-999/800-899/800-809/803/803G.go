package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n int
	t []int64
}

func NewSegTree(a []int64) *SegTree {
	n := len(a)
	st := &SegTree{n: n, t: make([]int64, 4*n)}
	var build func(id, l, r int)
	build = func(id, l, r int) {
		if l == r {
			st.t[id] = a[l-1]
			return
		}
		m := (l + r) >> 1
		build(id<<1, l, m)
		build(id<<1|1, m+1, r)
		if st.t[id<<1] < st.t[id<<1|1] {
			st.t[id] = st.t[id<<1]
		} else {
			st.t[id] = st.t[id<<1|1]
		}
	}
	build(1, 1, n)
	return st
}

func (st *SegTree) query(id, l, r, ql, qr int) int64 {
	if ql <= l && r <= qr {
		return st.t[id]
	}
	m := (l + r) >> 1
	res := int64(1<<63 - 1)
	if ql <= m {
		v := st.query(id<<1, l, m, ql, qr)
		if v < res {
			res = v
		}
	}
	if qr > m {
		v := st.query(id<<1|1, m+1, r, ql, qr)
		if v < res {
			res = v
		}
	}
	return res
}

func (st *SegTree) Query(l, r int) int64 {
	return st.query(1, 1, st.n, l, r)
}

type Node struct {
	left, right *Node
	lazy        bool
	val         int64
	min         int64
}

var (
	n, k int
	base *SegTree
	bMin int64
	n64  int64
	root *Node
)

func baseMinRange(L, R int64) int64 {
	lBlock := (L - 1) / n64
	rBlock := (R - 1) / n64
	lPos := int((L-1)%n64) + 1
	rPos := int((R-1)%n64) + 1
	if lBlock == rBlock {
		return base.Query(lPos, rPos)
	}
	res := base.Query(lPos, n)
	v := base.Query(1, rPos)
	if v < res {
		res = v
	}
	if rBlock > lBlock+1 {
		if bMin < res {
			res = bMin
		}
	}
	return res
}

func getMin(node *Node, l, r int64) int64 {
	if node == nil {
		return baseMinRange(l, r)
	}
	return node.min
}

func push(node *Node, l, r int64) {
	if node.lazy && l != r {
		if node.left == nil {
			node.left = &Node{}
		}
		if node.right == nil {
			node.right = &Node{}
		}
		node.left.lazy = true
		node.left.val = node.val
		node.left.min = node.val
		node.left.left = nil
		node.left.right = nil

		node.right.lazy = true
		node.right.val = node.val
		node.right.min = node.val
		node.right.left = nil
		node.right.right = nil

		node.lazy = false
	}
}

func update(node **Node, l, r, L, R int64, val int64) {
	if L > r || R < l {
		return
	}
	if *node == nil {
		*node = &Node{min: baseMinRange(l, r)}
	}
	if L <= l && r <= R {
		(*node).lazy = true
		(*node).val = val
		(*node).min = val
		(*node).left = nil
		(*node).right = nil
		return
	}
	push(*node, l, r)
	mid := (l + r) >> 1
	update(&(*node).left, l, mid, L, R, val)
	update(&(*node).right, mid+1, r, L, R, val)
	leftMin := getMin((*node).left, l, mid)
	rightMin := getMin((*node).right, mid+1, r)
	if leftMin < rightMin {
		(*node).min = leftMin
	} else {
		(*node).min = rightMin
	}
}

func query(node *Node, l, r, L, R int64) int64 {
	if L > r || R < l {
		return int64(1<<63 - 1)
	}
	if node == nil {
		ll := L
		if ll < l {
			ll = l
		}
		rr := R
		if rr > r {
			rr = r
		}
		return baseMinRange(ll, rr)
	}
	if L <= l && r <= R {
		return node.min
	}
	push(node, l, r)
	mid := (l + r) >> 1
	res := int64(1<<63 - 1)
	if L <= mid {
		v := query(node.left, l, mid, L, R)
		if v < res {
			res = v
		}
	}
	if R > mid {
		v := query(node.right, mid+1, r, L, R)
		if v < res {
			res = v
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n, &k)
	arr := make([]int64, n)
	bMin = int64(1<<63 - 1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] < bMin {
			bMin = arr[i]
		}
	}
	base = NewSegTree(arr)
	n64 = int64(n)
	total := n64 * int64(k)

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var l, r int64
			var x int64
			fmt.Fscan(in, &l, &r, &x)
			update(&root, 1, total, l, r, x)
		} else {
			var l, r int64
			fmt.Fscan(in, &l, &r)
			ans := query(root, 1, total, l, r)
			fmt.Fprintln(out, ans)
		}
	}
}

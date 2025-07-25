package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Node struct {
	key   int
	val   int64
	prio  int
	left  *Node
	right *Node
}

func rotateRight(p *Node) *Node {
	q := p.left
	p.left = q.right
	q.right = p
	return q
}

func rotateLeft(p *Node) *Node {
	q := p.right
	p.right = q.left
	q.left = p
	return q
}

func insertNode(root *Node, key int, val int64) *Node {
	if root == nil {
		return &Node{key: key, val: val, prio: rand.Int()}
	}
	if key < root.key {
		root.left = insertNode(root.left, key, val)
		if root.left.prio > root.prio {
			root = rotateRight(root)
		}
	} else {
		root.right = insertNode(root.right, key, val)
		if root.right.prio > root.prio {
			root = rotateLeft(root)
		}
	}
	return root
}

func search(root *Node, key int) *Node {
	for root != nil {
		if key < root.key {
			root = root.left
		} else if key > root.key {
			root = root.right
		} else {
			return root
		}
	}
	return nil
}

func predecessor(root *Node, key int) *Node {
	var res *Node
	for root != nil {
		if root.key < key {
			res = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *Node, key int) *Node {
	var res *Node
	for root != nil {
		if root.key > key {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

type SegTree struct {
	n      int
	sum    []int64
	lazyA  []int64
	lazyB  []int64
	prefix []int64
}

func newSegTree(n int) *SegTree {
	st := &SegTree{
		n:      n,
		sum:    make([]int64, 4*n),
		lazyA:  make([]int64, 4*n),
		lazyB:  make([]int64, 4*n),
		prefix: make([]int64, n+1),
	}
	for i := 1; i <= n; i++ {
		st.prefix[i] = st.prefix[i-1] + int64(i)
	}
	return st
}

func (st *SegTree) apply(idx, l, r int, A, B int64) {
	st.sum[idx] += A*(st.prefix[r]-st.prefix[l-1]) + B*int64(r-l+1)
	st.lazyA[idx] += A
	st.lazyB[idx] += B
}

func (st *SegTree) push(idx, l, r int) {
	if st.lazyA[idx] != 0 || st.lazyB[idx] != 0 {
		m := (l + r) >> 1
		st.apply(idx<<1, l, m, st.lazyA[idx], st.lazyB[idx])
		st.apply(idx<<1|1, m+1, r, st.lazyA[idx], st.lazyB[idx])
		st.lazyA[idx] = 0
		st.lazyB[idx] = 0
	}
}

func (st *SegTree) rangeAdd(idx, l, r, ql, qr int, A, B int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(idx, l, r, A, B)
		return
	}
	m := (l + r) >> 1
	st.push(idx, l, r)
	st.rangeAdd(idx<<1, l, m, ql, qr, A, B)
	st.rangeAdd(idx<<1|1, m+1, r, ql, qr, A, B)
	st.sum[idx] = st.sum[idx<<1] + st.sum[idx<<1|1]
}

func (st *SegTree) query(idx, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return st.sum[idx]
	}
	m := (l + r) >> 1
	st.push(idx, l, r)
	return st.query(idx<<1, l, m, ql, qr) + st.query(idx<<1|1, m+1, r, ql, qr)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	X := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &X[i])
	}
	V := make([]int64, m)
	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(in, &t)
		V[i] = int64(t)
	}

	// sort harbours by position
	idx := make([]int, m)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return X[idx[i]] < X[idx[j]] })

	var root *Node
	value := make(map[int]int64)
	pos := make([]int, m)
	for i, id := range idx {
		x := X[id]
		v := V[id]
		root = insertNode(root, x, v)
		pos[i] = x
		value[x] = v
	}

	st := newSegTree(n)
	// build initial costs
	for i := 0; i < len(pos)-1; i++ {
		L := pos[i]
		R := pos[i+1]
		v := value[L]
		if L+1 <= R-1 {
			st.rangeAdd(1, 1, n, L+1, R-1, -v, v*int64(R))
		}
	}

	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var x int
			var vv int64
			fmt.Fscan(in, &x, &vv)
			lnode := predecessor(root, x)
			rnode := successor(root, x)
			if lnode == nil || rnode == nil {
				continue
			}
			L := lnode.key
			R := rnode.key
			vl := lnode.val
			if L+1 <= x {
				st.rangeAdd(1, 1, n, L+1, x, 0, vl*int64(x-R))
			}
			if x+1 <= R-1 {
				delta := vv - vl
				st.rangeAdd(1, 1, n, x+1, R-1, -delta, delta*int64(R))
			}
			root = insertNode(root, x, vv)
			value[x] = vv
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			ans := st.query(1, 1, n, l, r)
			fmt.Fprintln(out, ans)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	v := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		v = v*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * v
}

func (fs *FastScanner) NextUint64() uint64 {
	var v uint64
	c, _ := fs.r.ReadByte()
	for c < '0' || c > '9' {
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		v = v*10 + uint64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return v
}

const maxOpsPerIndex = 40
const inf = int64(1 << 60)

type histEntry struct {
	node int
	prev int
}

type MinSegTree struct {
	size int
	tree []int64
}

func NewMinSegTree(capacity int) *MinSegTree {
	size := 1
	for size < capacity {
		size <<= 1
	}
	tree := make([]int64, size<<1)
	for i := range tree {
		tree[i] = inf
	}
	return &MinSegTree{
		size: size,
		tree: tree,
	}
}

func (st *MinSegTree) Update(pos int, val int64) {
	idx := pos + st.size
	st.tree[idx] = val
	idx >>= 1
	for idx > 0 {
		left := st.tree[idx<<1]
		right := st.tree[idx<<1|1]
		if left < right {
			st.tree[idx] = left
		} else {
			st.tree[idx] = right
		}
		idx >>= 1
	}
}

func (st *MinSegTree) Query(l, r int) int64 {
	if l > r {
		return inf
	}
	l += st.size
	r += st.size
	res := inf
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

type LiChao struct {
	xs      []uint64
	tree    []int
	history []histEntry
	a       []uint64
	base    *[]int64
}

func NewLiChao(xs []uint64, a []uint64, base *[]int64) *LiChao {
	m := len(xs)
	tree := make([]int, m*4)
	for i := range tree {
		tree[i] = -1
	}
	return &LiChao{
		xs:   xs,
		tree: tree,
		a:    a,
		base: base,
	}
}

type evalRes struct {
	intPart int64
	rem     uint64
	den     uint64
	cost    int64
}

func (lc *LiChao) eval(idx int, x uint64) evalRes {
	k := lc.a[idx-1]
	base := (*lc.base)[idx]
	div := x / k
	rem := x % k
	intPart := base + int64(div)
	cost := base + int64((x+k-1)/k)
	return evalRes{
		intPart: intPart,
		rem:     rem,
		den:     k,
		cost:    cost,
	}
}

func lessFrac(r1, d1, r2, d2 uint64) bool {
	h1, l1 := bits.Mul64(r1, d2)
	h2, l2 := bits.Mul64(r2, d1)
	if h1 != h2 {
		return h1 < h2
	}
	return l1 < l2
}

func (lc *LiChao) better(i, j int, x uint64) bool {
	if i == -1 {
		return false
	}
	if j == -1 {
		return true
	}
	v1 := lc.eval(i, x)
	v2 := lc.eval(j, x)
	if v1.intPart != v2.intPart {
		return v1.intPart < v2.intPart
	}
	if v1.rem == 0 && v2.rem == 0 {
		// equal
	} else if v1.rem == 0 {
		return true
	} else if v2.rem == 0 {
		return false
	} else if v1.den != 0 && v2.den != 0 {
		if v1.rem != v2.rem || v1.den != v2.den {
			return lessFrac(v1.rem, v1.den, v2.rem, v2.den)
		}
	}
	if v1.cost != v2.cost {
		return v1.cost < v2.cost
	}
	return i < j
}

func (lc *LiChao) snapshot() int {
	return len(lc.history)
}

func (lc *LiChao) rollback(sz int) {
	for len(lc.history) > sz {
		last := lc.history[len(lc.history)-1]
		lc.history = lc.history[:len(lc.history)-1]
		lc.tree[last.node] = last.prev
	}
}

func (lc *LiChao) setNode(node, val int) {
	lc.history = append(lc.history, histEntry{node: node, prev: lc.tree[node]})
	lc.tree[node] = val
}

func (lc *LiChao) insertLine(idx int) {
	if len(lc.xs) == 0 {
		return
	}
	lc.insert(1, 0, len(lc.xs)-1, idx)
}

func (lc *LiChao) insert(node, l, r, idx int) {
	cur := lc.tree[node]
	if cur == -1 {
		lc.setNode(node, idx)
		return
	}
	mid := (l + r) >> 1
	xMid := lc.xs[mid]
	if lc.better(idx, cur, xMid) {
		lc.setNode(node, idx)
		idx = cur
	}
	if l == r {
		return
	}
	left := node << 1
	right := left | 1
	if lc.better(idx, lc.tree[node], lc.xs[l]) {
		lc.insert(left, l, mid, idx)
	} else if lc.better(idx, lc.tree[node], lc.xs[r]) {
		lc.insert(right, mid+1, r, idx)
	}
}

func (lc *LiChao) query(pos int) int64 {
	x := lc.xs[pos]
	res := inf
	node := 1
	l, r := 0, len(lc.xs)-1
	for {
		cur := lc.tree[node]
		if cur != -1 {
			val := lc.eval(cur, x).cost
			if val < res {
				res = val
			}
		}
		if l == r {
			break
		}
		mid := (l + r) >> 1
		if pos <= mid {
			node = node << 1
			r = mid
		} else {
			node = node<<1 | 1
			l = mid + 1
		}
	}
	return res
}

func uniqueSorted(vals []uint64) []uint64 {
	if len(vals) == 0 {
		return vals
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	k := 1
	for i := 1; i < len(vals); i++ {
		if vals[i] != vals[k-1] {
			vals[k] = vals[i]
			k++
		}
	}
	return vals[:k]
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()
	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		a := make([]uint64, n)
		xs := make([]uint64, n)
		for i := 0; i < n; i++ {
			val := fs.NextUint64()
			a[i] = val
			xs[i] = val
		}
		nextLess := make([]int, n)
		stack := make([]int, 0)
		for i := n - 1; i >= 0; i-- {
			for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				nextLess[i] = n + 1
			} else {
				nextLess[i] = stack[len(stack)-1] + 1
			}
			stack = append(stack, i)
		}
		prevLess := make([]int, n+1)
		stack = stack[:0]
		for i := 1; i <= n; i++ {
			for len(stack) > 0 && a[stack[len(stack)-1]-1] >= a[i-1] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				prevLess[i] = 0
			} else {
				prevLess[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		xs = uniqueSorted(xs)
		posIdx := make([]int, n)
		posMap := make(map[uint64]int, len(xs))
		for i, v := range xs {
			posMap[v] = i
		}
		for i := 0; i < n; i++ {
			posIdx[i] = posMap[a[i]]
		}
		segSize := 4*n + 5
		count := 0
		var countRange func(node, l, r, L, R int)
		countRange = func(node, l, r, L, R int) {
			if L > R {
				return
			}
			if L <= l && r <= R {
				count++
				return
			}
			mid := (l + r) >> 1
			if L <= mid {
				countRange(node<<1, l, mid, L, R)
			}
			if R > mid {
				countRange(node<<1|1, mid+1, r, L, R)
			}
		}
		for j := 1; j <= n; j++ {
			L := j
			R := nextLess[j-1] - 1
			if L <= R {
				countRange(1, 1, n, L, R)
			}
		}
		head := make([]int32, segSize)
		for i := range head {
			head[i] = -1
		}
		to := make([]int32, count)
		nxt := make([]int32, count)
		ptr := 0
		var addRange func(node, l, r, L, R, idx int)
		addRange = func(node, l, r, L, R, idx int) {
			if L > R {
				return
			}
			if L <= l && r <= R {
				to[ptr] = int32(idx)
				nxt[ptr] = head[node]
				head[node] = int32(ptr)
				ptr++
				return
			}
			mid := (l + r) >> 1
			if L <= mid {
				addRange(node<<1, l, mid, L, R, idx)
			}
			if R > mid {
				addRange(node<<1|1, mid+1, r, L, R, idx)
			}
		}
		for j := 1; j <= n; j++ {
			L := j
			R := nextLess[j-1] - 1
			if L <= R {
				addRange(1, 1, n, L, R, j)
			}
		}

		dp := make([]int64, n+1)
		dp[0] = 0
		lineBase := make([]int64, n+1)
		for i := range lineBase {
			lineBase[i] = inf
		}
		li := NewLiChao(xs, a, &lineBase)
		dpSeg := NewMinSegTree(n + 1)
		dpSeg.Update(0, 0)
		var dfs func(node, l, r int)
		dfs = func(node, l, r int) {
			snap := li.snapshot()
			for cur := head[node]; cur != -1; {
				ci := int(cur)
				idx := int(to[ci])
				if lineBase[idx] == inf {
					best := dpSeg.Query(prevLess[idx], idx-1)
					lineBase[idx] = best
				}
				li.insertLine(idx)
				cur = nxt[ci]
			}
			if l == r {
				best := li.query(posIdx[l-1])
				dp[l] = best
				dpSeg.Update(l, dp[l])
			} else {
				mid := (l + r) >> 1
				dfs(node<<1, l, mid)
				dfs(node<<1|1, mid+1, r)
			}
			li.rollback(snap)
		}
		dfs(1, 1, n)
		fmt.Fprintln(out, dp[n])
	}
}

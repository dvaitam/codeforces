package main

import (
	"bufio"
	"fmt"
	"os"
)

type ordNode struct {
	key      int
	priority uint32
	left     *ordNode
	right    *ordNode
}

type multiNode struct {
	key      int
	priority uint32
	cnt      int
	size     int
	left     *multiNode
	right    *multiNode
}

var rng uint32 = 123456789

func nextRand() uint32 {
	rng ^= rng << 13
	rng ^= rng >> 17
	rng ^= rng << 5
	return rng
}

func ordRotateRight(t *ordNode) *ordNode {
	l := t.left
	t.left = l.right
	l.right = t
	return l
}

func ordRotateLeft(t *ordNode) *ordNode {
	r := t.right
	t.right = r.left
	r.left = t
	return r
}

func ordInsert(t *ordNode, key int) *ordNode {
	if t == nil {
		return &ordNode{key: key, priority: nextRand()}
	}
	if key < t.key {
		t.left = ordInsert(t.left, key)
		if t.left.priority < t.priority {
			t = ordRotateRight(t)
		}
	} else if key > t.key {
		t.right = ordInsert(t.right, key)
		if t.right.priority < t.priority {
			t = ordRotateLeft(t)
		}
	}
	return t
}

func ordErase(t *ordNode, key int) *ordNode {
	if t == nil {
		return nil
	}
	if key < t.key {
		t.left = ordErase(t.left, key)
	} else if key > t.key {
		t.right = ordErase(t.right, key)
	} else {
		if t.left == nil {
			return t.right
		}
		if t.right == nil {
			return t.left
		}
		if t.left.priority < t.right.priority {
			t = ordRotateRight(t)
			t.right = ordErase(t.right, key)
		} else {
			t = ordRotateLeft(t)
			t.left = ordErase(t.left, key)
		}
	}
	return t
}

func ordPredecessor(t *ordNode, key int) int {
	res := 0
	for t != nil {
		if key <= t.key {
			t = t.left
		} else {
			res = t.key
			t = t.right
		}
	}
	return res
}

func ordSuccessor(t *ordNode, key int) int {
	res := 0
	for t != nil {
		if key >= t.key {
			t = t.right
		} else {
			res = t.key
			t = t.left
		}
	}
	return res
}

func multiRotateRight(t *multiNode) *multiNode {
	l := t.left
	t.left = l.right
	l.right = t
	updateSize(t)
	updateSize(l)
	return l
}

func multiRotateLeft(t *multiNode) *multiNode {
	r := t.right
	t.right = r.left
	r.left = t
	updateSize(t)
	updateSize(r)
	return r
}

func updateSize(t *multiNode) {
	if t != nil {
		t.size = t.cnt + getSize(t.left) + getSize(t.right)
	}
}

func getSize(t *multiNode) int {
	if t == nil {
		return 0
	}
	return t.size
}

func multiInsert(t *multiNode, key int) *multiNode {
	if t == nil {
		return &multiNode{key: key, priority: nextRand(), cnt: 1, size: 1}
	}
	if key < t.key {
		t.left = multiInsert(t.left, key)
		if t.left.priority < t.priority {
			t = multiRotateRight(t)
		}
	} else if key > t.key {
		t.right = multiInsert(t.right, key)
		if t.right.priority < t.priority {
			t = multiRotateLeft(t)
		}
	} else {
		t.cnt++
	}
	updateSize(t)
	return t
}

func multiErase(t *multiNode, key int) *multiNode {
	if t == nil {
		return nil
	}
	if key < t.key {
		t.left = multiErase(t.left, key)
	} else if key > t.key {
		t.right = multiErase(t.right, key)
	} else {
		if t.cnt > 1 {
			t.cnt--
		} else {
			if t.left == nil {
				return t.right
			}
			if t.right == nil {
				return t.left
			}
			if t.left.priority < t.right.priority {
				t = multiRotateRight(t)
				t.right = multiErase(t.right, key)
			} else {
				t = multiRotateLeft(t)
				t.left = multiErase(t.left, key)
			}
		}
	}
	updateSize(t)
	return t
}

func multiCountLess(t *multiNode, key int) int {
	if t == nil {
		return 0
	}
	if key <= t.key {
		return multiCountLess(t.left, key)
	}
	return getSize(t.left) + t.cnt + multiCountLess(t.right, key)
}

type segTree struct {
	n    int
	tree []*multiNode
}

func newSegTree(n int, values []int) *segTree {
	st := &segTree{n: n, tree: make([]*multiNode, 4*n)}
	for i := 1; i <= n; i++ {
		st.add(1, 1, n, i, values[i])
	}
	return st
}

func (st *segTree) add(node, l, r, pos, val int) {
	st.tree[node] = multiInsert(st.tree[node], val)
	if l == r {
		return
	}
	mid := (l + r) >> 1
	if pos <= mid {
		st.add(node<<1, l, mid, pos, val)
	} else {
		st.add(node<<1|1, mid+1, r, pos, val)
	}
}

func (st *segTree) replace(node, l, r, pos, oldVal, newVal int) {
	st.tree[node] = multiErase(st.tree[node], oldVal)
	st.tree[node] = multiInsert(st.tree[node], newVal)
	if l == r {
		return
	}
	mid := (l + r) >> 1
	if pos <= mid {
		st.replace(node<<1, l, mid, pos, oldVal, newVal)
	} else {
		st.replace(node<<1|1, mid+1, r, pos, oldVal, newVal)
	}
}

func (st *segTree) countLess(node, l, r, ql, qr, threshold int) int {
	if ql <= l && r <= qr {
		return multiCountLess(st.tree[node], threshold)
	}
	mid := (l + r) >> 1
	res := 0
	if ql <= mid {
		res += st.countLess(node<<1, l, mid, ql, qr, threshold)
	}
	if qr > mid {
		res += st.countLess(node<<1|1, mid+1, r, ql, qr, threshold)
	}
	return res
}

type solver struct {
	n     int
	prev  []int
	a     []int
	roots []*ordNode
	seg   *segTree
}

func newSolver(n int, arr []int) *solver {
	s := &solver{
		n:     n,
		prev:  make([]int, n+1),
		a:     make([]int, n+1),
		roots: make([]*ordNode, n+1),
	}
	copy(s.a, arr)
	last := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v := s.a[i]
		s.prev[i] = last[v]
		last[v] = i
		s.roots[v] = ordInsert(s.roots[v], i)
	}
	s.seg = newSegTree(n, s.prev)
	return s
}

func (s *solver) setPrev(pos, value int) {
	if pos == 0 {
		return
	}
	old := s.prev[pos]
	if old == value {
		return
	}
	s.seg.replace(1, 1, s.n, pos, old, value)
	s.prev[pos] = value
}

func (s *solver) countPrevLess(l, r, threshold int) int {
	return s.seg.countLess(1, 1, s.n, l, r, threshold)
}

func (s *solver) changeValue(pos, newVal int) {
	oldVal := s.a[pos]
	if oldVal == newVal {
		return
	}
	rootOld := s.roots[oldVal]
	preOld := ordPredecessor(rootOld, pos)
	sucOld := ordSuccessor(rootOld, pos)
	s.roots[oldVal] = ordErase(rootOld, pos)
	if sucOld != 0 {
		s.setPrev(sucOld, preOld)
	}

	rootNew := s.roots[newVal]
	preNew := ordPredecessor(rootNew, pos)
	sucNew := ordSuccessor(rootNew, pos)
	s.roots[newVal] = ordInsert(rootNew, pos)
	s.setPrev(pos, preNew)
	if sucNew != 0 {
		s.setPrev(sucNew, pos)
	}
	s.a[pos] = newVal
}

func nextInt(r *bufio.Reader) int {
	sign := 1
	val := 0
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = r.ReadByte()
	}
	return sign * val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := nextInt(in)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = nextInt(in)
	}
	s := newSolver(n, arr)
	q := nextInt(in)
	last := int64(0)
	modN := int64(n)
	for i := 0; i < q; i++ {
		typ := nextInt(in)
		if typ == 1 {
			pPrime := nextInt(in)
			xPrime := nextInt(in)
			p := int((int64(pPrime) + last) % modN)
			x := int((int64(xPrime) + last) % modN)
			p++
			x++
			s.changeValue(p, x)
		} else {
			lPrime := nextInt(in)
			rPrime := nextInt(in)
			l := int((int64(lPrime) + last) % modN)
			r := int((int64(rPrime) + last) % modN)
			l++
			r++
			if l > r {
				l, r = r, l
			}
			length := r - l + 1
			cntLess := s.countPrevLess(l, r, l)
			equalPairs := length - cntLess
			totalPairs := int64(length) * int64(length-1) / 2
			ans := totalPairs - int64(equalPairs)
			fmt.Fprintln(out, ans)
			last = ans
		}
	}
}

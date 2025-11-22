package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type Info struct {
	size     int
	maxVal   int64
	minVal   int64
	maxPref  int64
	minSuf   int64
	bestDiff int64 // max over i<=j of val_i - val_j
}

func leafInfo(val int64) Info {
	return Info{
		size:     1,
		maxVal:   val,
		minVal:   val,
		maxPref:  val,
		minSuf:   val,
		bestDiff: 0,
	}
}

func shiftInfo(a Info, delta int64) Info {
	if a.size == 0 {
		return a
	}
	a.maxVal += delta
	a.minVal += delta
	a.maxPref += delta
	a.minSuf += delta
	return a
}

func concatInfo(a, b Info) Info {
	if a.size == 0 {
		return b
	}
	if b.size == 0 {
		return a
	}
	res := Info{size: a.size + b.size}
	if a.maxVal > b.maxVal {
		res.maxVal = a.maxVal
	} else {
		res.maxVal = b.maxVal
	}
	if a.minVal < b.minVal {
		res.minVal = a.minVal
	} else {
		res.minVal = b.minVal
	}
	// max prefix can end in a, at boundary, or inside b
	maxValBoundary := a.maxVal
	if b.maxPref > maxValBoundary {
		maxValBoundary = b.maxPref
	}
	if a.maxPref > maxValBoundary {
		res.maxPref = a.maxPref
	} else {
		res.maxPref = maxValBoundary
	}
	// min suffix can start in b, at boundary, or inside a
	minValBoundary := b.minVal
	if a.minSuf < minValBoundary {
		minValBoundary = a.minSuf
	}
	if b.minSuf < minValBoundary {
		res.minSuf = b.minSuf
	} else {
		res.minSuf = minValBoundary
	}

	// best diff: inside a, inside b, or crossing
	cross := a.maxPref - b.minSuf
	res.bestDiff = a.bestDiff
	if b.bestDiff > res.bestDiff {
		res.bestDiff = b.bestDiff
	}
	if cross > res.bestDiff {
		res.bestDiff = cross
	}
	return res
}

type Node struct {
	pos       int
	priority  uint32
	left      *Node
	right     *Node
	info      Info
	subtreeSz int
}

func getInfo(t *Node) Info {
	if t == nil {
		return Info{}
	}
	return t.info
}

func pull(t *Node) {
	if t == nil {
		return
	}
	lInfo := getInfo(t.left)
	rInfo := getInfo(t.right)
	// compute value for root with correct index shift
	valRoot := int64(t.pos) - 2*int64(lInfo.size+1)
	rootInfo := leafInfo(valRoot)
	shiftedRight := shiftInfo(rInfo, -2*int64(lInfo.size+1))

	tmp := concatInfo(lInfo, rootInfo)
	t.info = concatInfo(tmp, shiftedRight)
	t.subtreeSz = lInfo.size + 1 + rInfo.size
}

// split by key, left contains positions < key
func split(t *Node, key int) (*Node, *Node) {
	if t == nil {
		return nil, nil
	}
	if key <= t.pos {
		l, r := split(t.left, key)
		t.left = r
		pull(t)
		return l, t
	}
	l, r := split(t.right, key)
	t.right = l
	pull(t)
	return t, r
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		a.right = merge(a.right, b)
		pull(a)
		return a
	}
	b.left = merge(a, b.left)
	pull(b)
	return b
}

func insert(t *Node, pos int) *Node {
	node := &Node{pos: pos, priority: rand.Uint32()}
	left, right := split(t, pos)
	return merge(merge(left, node), right)
}

func erase(t *Node, pos int) *Node {
	if t == nil {
		return nil
	}
	if pos == t.pos {
		return merge(t.left, t.right)
	}
	if pos < t.pos {
		t.left = erase(t.left, pos)
	} else {
		t.right = erase(t.right, pos)
	}
	pull(t)
	return t
}

func kVal(root *Node) int {
	if root == nil || root.info.size == 0 {
		return 0
	}
	return int((root.info.bestDiff + 1) / 2)
}

// segment tree for max
type SegTree struct {
	n int
	t []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: size, t: make([]int, size<<1)}
}

func (s *SegTree) Build(arr []int) {
	for i := 0; i < len(arr); i++ {
		s.t[s.n+i] = arr[i]
	}
	for i := s.n - 1; i > 0; i-- {
		if s.t[i<<1] > s.t[i<<1|1] {
			s.t[i] = s.t[i<<1]
		} else {
			s.t[i] = s.t[i<<1|1]
		}
	}
}

func (s *SegTree) Update(pos, val int) {
	p := s.n + pos
	s.t[p] = val
	for p >>= 1; p > 0; p >>= 1 {
		left := s.t[p<<1]
		right := s.t[p<<1|1]
		if left > right {
			s.t[p] = left
		} else {
			s.t[p] = right
		}
	}
}

func (s *SegTree) Max() int {
	return s.t[1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		roots := make([]*Node, n+2)
		for i, v := range a {
			roots[v] = insert(roots[v], i+1)
		}

		kArr := make([]int, n+2)
		for v := 1; v <= n; v++ {
			kArr[v] = kVal(roots[v])
		}
		seg := NewSegTree(n + 2)
		seg.Build(kArr)

		for qi := 0; qi < q; qi++ {
			var idx, x int
			fmt.Fscan(in, &idx, &x)
			old := a[idx-1]
			if old != x {
				roots[old] = erase(roots[old], idx)
				seg.Update(old, kVal(roots[old]))
				roots[x] = insert(roots[x], idx)
				seg.Update(x, kVal(roots[x]))
				a[idx-1] = x
			}
			if qi > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, seg.Max())
		}
		fmt.Fprintln(out)
	}
}

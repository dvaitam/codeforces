package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	var val int64
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

type Interval struct {
	l, r int64
	idx  int
}

type Channel struct {
	a, b, c int64
	idx     int
}

type Pair struct {
	val int64
	idx int
}

func better(a, b Pair) Pair {
	if a.val > b.val {
		return a
	}
	return b
}

type Fenwick struct {
	n   int
	bit []Pair
}

func NewFenwick(n int) *Fenwick {
	bit := make([]Pair, n+2)
	for i := range bit {
		bit[i].val = -1
		bit[i].idx = -1
	}
	return &Fenwick{n: n + 2, bit: bit}
}

func (f *Fenwick) update(pos int, val int64, idx int) {
	if val < 0 {
		return
	}
	for pos < f.n {
		if val > f.bit[pos].val {
			f.bit[pos].val = val
			f.bit[pos].idx = idx
		}
		pos += pos & -pos
	}
}

func (f *Fenwick) query(pos int) Pair {
	res := Pair{-1, -1}
	for pos > 0 {
		if f.bit[pos].val > res.val {
			res = f.bit[pos]
		}
		pos -= pos & -pos
	}
	return res
}

// Segment tree for case1 (max value in range)

type SegTreeVals struct {
	n    int
	tree []Pair
}

func NewSegTreeVals(n int) *SegTreeVals {
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]Pair, 2*size)
	for i := range tree {
		tree[i].val = -1
		tree[i].idx = -1
	}
	return &SegTreeVals{n: size, tree: tree}
}

func (st *SegTreeVals) update(pos int, val int64, idx int) {
	if val < 0 {
		return
	}
	pos += st.n
	if val > st.tree[pos].val {
		st.tree[pos] = Pair{val, idx}
	}
	pos >>= 1
	for pos > 0 {
		left := st.tree[pos<<1]
		right := st.tree[pos<<1|1]
		if left.val >= right.val {
			st.tree[pos] = left
		} else {
			st.tree[pos] = right
		}
		pos >>= 1
	}
}

func (st *SegTreeVals) query(l, r int) Pair {
	res := Pair{-1, -1}
	if l > r {
		return res
	}
	l += st.n
	r += st.n
	for l <= r {
		if l&1 == 1 {
			if st.tree[l].val > res.val {
				res = st.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.tree[r].val > res.val {
				res = st.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

type SegTreeMax struct {
	size int
	max  []Pair // val = max r, idx = index in l-order
}

func NewSegTreeMax(values []int64, idxs []int) *SegTreeMax {
	n := len(values)
	size := 1
	for size < n {
		size <<= 1
	}
	maxArr := make([]Pair, 2*size)
	for i := range maxArr {
		maxArr[i].val = -1
		maxArr[i].idx = -1
	}
	st := &SegTreeMax{size: size, max: maxArr}
	for i := 0; i < n; i++ {
		st.max[size+i] = Pair{values[i], idxs[i]}
	}
	for i := size - 1; i >= 1; i-- {
		left := st.max[i<<1]
		right := st.max[i<<1|1]
		if left.val >= right.val {
			st.max[i] = left
		} else {
			st.max[i] = right
		}
	}
	return st
}

func (st *SegTreeMax) queryMax(l, r int) int64 {
	if l > r {
		return -1
	}
	l += st.size
	r += st.size
	res := int64(-1)
	for l <= r {
		if l&1 == 1 {
			if st.max[l].val > res {
				res = st.max[l].val
			}
			l++
		}
		if r&1 == 0 {
			if st.max[r].val > res {
				res = st.max[r].val
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func (st *SegTreeMax) findFirstGreater(l, r int, threshold int64) int {
	return st._find(1, 0, st.size-1, l, r, threshold)
}

func (st *SegTreeMax) _find(node, nl, nr, ql, qr int, threshold int64) int {
	if ql > nr || qr < nl || st.max[node].val <= threshold {
		return -1
	}
	if nl == nr {
		return nl
	}
	mid := (nl + nr) >> 1
	left := st._find(node<<1, nl, mid, ql, qr, threshold)
	if left != -1 {
		return left
	}
	return st._find(node<<1|1, mid+1, nr, ql, qr, threshold)
}

func lowerBound(arr []int64, target int64) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] >= target })
}

func upperBound(arr []int64, target int64) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] > target })
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := int(in.NextInt64())
	m := int(in.NextInt64())

	intervals := make([]Interval, n)
	lVals := make([]int64, n)
	for i := 0; i < n; i++ {
		intervals[i] = Interval{l: in.NextInt64(), r: in.NextInt64(), idx: i + 1}
		lVals[i] = intervals[i].l
	}

	channels := make([]Channel, m)
	for i := 0; i < m; i++ {
		channels[i] = Channel{a: in.NextInt64(), b: in.NextInt64(), c: in.NextInt64(), idx: i + 1}
	}

	sort.Slice(lVals, func(i, j int) bool { return lVals[i] < lVals[j] })
	lVals = unique(lVals)
	coord := func(x int64) int {
		return lowerBound(lVals, x)
	}

	// case3 & 4 preparations: intervals sorted by l
	byL := make([]int, n)
	for i := range byL {
		byL[i] = i
	}
	sort.Slice(byL, func(i, j int) bool {
		li := intervals[byL[i]].l
		lj := intervals[byL[j]].l
		if li == lj {
			return intervals[byL[i]].r < intervals[byL[j]].r
		}
		return li < lj
	})
	lSorted := make([]int64, n)
	rSorted := make([]int64, n)
	idxSorted := make([]int, n)
	prefMaxR := make([]int64, n)
	prefIdx := make([]int, n)
	for i := 0; i < n; i++ {
		interval := intervals[byL[i]]
		lSorted[i] = interval.l
		rSorted[i] = interval.r
		idxSorted[i] = interval.idx
		prefMaxR[i] = interval.r
		prefIdx[i] = interval.idx
		if i > 0 && prefMaxR[i-1] >= prefMaxR[i] {
			prefMaxR[i] = prefMaxR[i-1]
			prefIdx[i] = prefIdx[i-1]
		}
	}
	segCase3 := NewSegTreeMax(rSorted, idxSorted)

	// case1 & case2 structures
	segCase1 := NewSegTreeVals(len(lVals))
	fen := NewFenwick(len(lVals) + 2)
	byR := make([]int, n)
	for i := range byR {
		byR[i] = i
	}
	sort.Slice(byR, func(i, j int) bool {
		ri := intervals[byR[i]].r
		rj := intervals[byR[j]].r
		if ri == rj {
			return intervals[byR[i]].l < intervals[byR[j]].l
		}
		return ri < rj
	})

	sortedChannels := make([]int, m)
	for i := range sortedChannels {
		sortedChannels[i] = i
	}
	sort.Slice(sortedChannels, func(i, j int) bool {
		bi := channels[sortedChannels[i]].b
		bj := channels[sortedChannels[j]].b
		if bi == bj {
			return channels[sortedChannels[i]].a < channels[sortedChannels[j]].a
		}
		return bi < bj
	})

	ptr := 0
	bestVal := int64(0)
	bestVideo := -1
	bestChannel := -1

	for _, idxCh := range sortedChannels {
		ch := channels[idxCh]
		for ptr < n && intervals[byR[ptr]].r <= ch.b {
			interval := intervals[byR[ptr]]
			pos := coord(interval.l)
			segCase1.update(pos, interval.r-interval.l, interval.idx)
			fen.update(pos+1, interval.r, interval.idx)
			ptr++
		}

		// case1
		posA := coord(ch.a)
		if posA < len(lVals) {
			res := segCase1.query(posA, len(lVals)-1)
			length := res.val
			if length > 0 {
				eff := length * ch.c
				if eff > bestVal {
					bestVal = eff
					bestVideo = res.idx
					bestChannel = ch.idx
				}
			}
		}
		// case2
		posLess := coord(ch.a)
		if posLess > 0 {
			res := fen.query(posLess)
			if res.val > ch.a {
				length := res.val - ch.a
				eff := length * ch.c
				if eff > bestVal {
					bestVal = eff
					bestVideo = res.idx
					bestChannel = ch.idx
				}
			}
		}

		// case3
		left := lowerBound(lSorted, ch.a)
		right := lowerBound(lSorted, ch.b)
		if left < right {
			maxVal := segCase3.queryMax(left, right-1)
			if maxVal > ch.b {
				idxPos := segCase3.findFirstGreater(left, right-1, ch.b)
				if idxPos != -1 {
					lVal := lSorted[idxPos]
					if lVal < ch.b {
						length := ch.b - lVal
						if length > 0 {
							eff := length * ch.c
							if eff > bestVal {
								bestVal = eff
								bestVideo = idxSorted[idxPos]
								bestChannel = ch.idx
							}
						}
					}
				}
			}
		}

		// case4
		posCover := lowerBound(lSorted, ch.a)
		if posCover > 0 && ch.b > ch.a {
			if prefMaxR[posCover-1] > ch.b {
				length := ch.b - ch.a
				eff := length * ch.c
				if eff > bestVal {
					bestVal = eff
					bestVideo = prefIdx[posCover-1]
					bestChannel = ch.idx
				}
			}
		}
	}

	if bestVal <= 0 {
		fmt.Fprintln(out, 0)
		return
	}
	fmt.Fprintln(out, bestVal)
	fmt.Fprintf(out, "%d %d\n", bestVideo, bestChannel)
}

func unique(arr []int64) []int64 {
	if len(arr) == 0 {
		return arr
	}
	res := arr[:1]
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			res = append(res, arr[i])
		}
	}
	return res
}

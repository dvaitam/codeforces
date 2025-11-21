package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner(reader io.Reader) *FastScanner {
	return &FastScanner{r: bufio.NewReader(reader)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

type Candidate struct {
	val int
	cnt int
}

type Node struct {
	c [2]Candidate
}

func (n *Node) add(val, cnt int) {
	if cnt == 0 {
		return
	}
	for i := 0; i < 2; i++ {
		if n.c[i].cnt > 0 && n.c[i].val == val {
			n.c[i].cnt += cnt
			return
		}
	}
	for i := 0; i < 2; i++ {
		if n.c[i].cnt == 0 {
			n.c[i] = Candidate{val: val, cnt: cnt}
			return
		}
	}
	minCount := n.c[0].cnt
	if n.c[1].cnt < minCount {
		minCount = n.c[1].cnt
	}
	if cnt <= minCount {
		n.c[0].cnt -= cnt
		n.c[1].cnt -= cnt
		return
	}
	n.c[0].cnt -= minCount
	n.c[1].cnt -= minCount
	cnt -= minCount
	for i := 0; i < 2; i++ {
		if n.c[i].cnt == 0 {
			n.c[i] = Candidate{val: val, cnt: cnt}
			return
		}
	}
	n.c[0] = Candidate{val: val, cnt: cnt}
}

func mergeNodes(a, b Node) Node {
	res := a
	for i := 0; i < 2; i++ {
		if b.c[i].cnt > 0 {
			res.add(b.c[i].val, b.c[i].cnt)
		}
	}
	return res
}

type SegmentTree struct {
	n    int
	tree []Node
	arr  []int
}

func NewSegmentTree(arr []int) *SegmentTree {
	st := &SegmentTree{
		n:    len(arr),
		tree: make([]Node, len(arr)*4),
		arr:  arr,
	}
	if st.n > 0 {
		st.build(1, 0, st.n-1)
	}
	return st
}

func (st *SegmentTree) build(idx, l, r int) {
	if l == r {
		st.tree[idx].add(st.arr[l], 1)
		return
	}
	mid := (l + r) >> 1
	st.build(idx<<1, l, mid)
	st.build(idx<<1|1, mid+1, r)
	st.tree[idx] = mergeNodes(st.tree[idx<<1], st.tree[idx<<1|1])
}

func (st *SegmentTree) query(ql, qr int) Node {
	if ql < 0 {
		ql = 0
	}
	if qr >= st.n {
		qr = st.n - 1
	}
	return st.queryRec(1, 0, st.n-1, ql, qr)
}

func (st *SegmentTree) queryRec(idx, l, r, ql, qr int) Node {
	if ql <= l && r <= qr {
		return st.tree[idx]
	}
	mid := (l + r) >> 1
	if qr <= mid {
		return st.queryRec(idx<<1, l, mid, ql, qr)
	}
	if ql > mid {
		return st.queryRec(idx<<1|1, mid+1, r, ql, qr)
	}
	left := st.queryRec(idx<<1, l, mid, ql, qr)
	right := st.queryRec(idx<<1|1, mid+1, r, ql, qr)
	return mergeNodes(left, right)
}

func countInRange(pos []int, l, r int) int {
	left := sort.Search(len(pos), func(i int) bool { return pos[i] >= l })
	right := sort.Search(len(pos), func(i int) bool { return pos[i] > r })
	return right - left
}

func main() {
	scanner := NewFastScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := scanner.NextInt()
	for ; t > 0; t-- {
		n := scanner.NextInt()
		q := scanner.NextInt()
		arr := make([]int, n)
		occ := make(map[int][]int)
		for i := 0; i < n; i++ {
			arr[i] = scanner.NextInt()
			occ[arr[i]] = append(occ[arr[i]], i+1)
		}
		st := NewSegmentTree(arr)
		for ; q > 0; q-- {
			l := scanner.NextInt()
			r := scanner.NextInt()
			node := st.query(l-1, r-1)
			length := r - l + 1
			res := make([]int, 0, 2)
			for i := 0; i < 2; i++ {
				cand := node.c[i]
				if cand.cnt == 0 {
					continue
				}
				val := cand.val
				freq := countInRange(occ[val], l, r)
				if freq*3 > length {
					exists := false
					for _, v := range res {
						if v == val {
							exists = true
							break
						}
					}
					if !exists {
						res = append(res, val)
					}
				}
			}
			if len(res) == 0 {
				fmt.Fprintln(out, -1)
			} else {
				sort.Ints(res)
				for i := 0; i < len(res); i++ {
					if i > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, res[i])
				}
				fmt.Fprintln(out)
			}
		}
	}
}

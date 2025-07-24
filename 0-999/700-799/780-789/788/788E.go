package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

// Fenwick tree for prefix sums of integers
type BIT struct {
	n   int
	bit []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, bit: make([]int, n+2)}
}

func (b *BIT) Add(idx, delta int) {
	idx++
	for idx <= b.n+1 {
		b.bit[idx] += delta
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	if idx < 0 {
		return 0
	}
	s := 0
	idx++
	if idx > b.n+1 {
		idx = b.n + 1
	}
	for idx > 0 {
		s += b.bit[idx]
		idx -= idx & -idx
	}
	return s
}

type Node struct {
	cnt  int64
	sumL int64
	sumR int64
	pref int64
	suff int64
	tot  int64
}

func merge(a, b Node) Node {
	var res Node
	res.cnt = a.cnt + b.cnt
	res.sumL = a.sumL + b.sumL
	res.sumR = a.sumR + b.sumR
	res.pref = a.pref + b.pref + b.cnt*a.sumL
	res.suff = a.suff + b.suff + a.cnt*b.sumR
	res.tot = a.tot + b.tot + b.sumR*a.pref + a.sumL*b.suff
	return res
}

type SegTree struct {
	size int
	tree []Node
	L    []int64
	R    []int64
}

func BuildSegTree(L, R []int64) *SegTree {
	n := len(L)
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]Node, 2*size)
	st := &SegTree{size: size, tree: tree, L: L, R: R}
	for i := 0; i < n; i++ {
		v := size + i
		tree[v] = Node{cnt: 1, sumL: L[i], sumR: R[i]}
	}
	for v := size - 1; v > 0; v-- {
		tree[v] = merge(tree[2*v], tree[2*v+1])
	}
	return st
}

func (st *SegTree) Update(idx int, active bool) {
	v := st.size + idx
	if active {
		st.tree[v] = Node{cnt: 1, sumL: st.L[idx], sumR: st.R[idx]}
	} else {
		st.tree[v] = Node{}
	}
	for v >>= 1; v > 0; v >>= 1 {
		st.tree[v] = merge(st.tree[2*v], st.tree[2*v+1])
	}
}

func (st *SegTree) Total() int64 {
	if len(st.tree) == 0 {
		return 0
	}
	return st.tree[1].tot
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	skills := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &skills[i])
	}

	// Compress skills
	uniq := append([]int(nil), skills...)
	sort.Ints(uniq)
	uniq = uniqueInts(uniq)
	mSkill := make(map[int]int)
	for i, v := range uniq {
		mSkill[v] = i
	}
	S := len(uniq)

	// Precompute left and right counts
	bit := NewBIT(S)
	L := make([]int64, n)
	for i := 0; i < n; i++ {
		s := mSkill[skills[i]]
		L[i] = int64(bit.Sum(s))
		bit.Add(s, 1)
	}

	bit = NewBIT(S)
	R := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		s := mSkill[skills[i]]
		R[i] = int64(bit.Sum(s))
		bit.Add(s, 1)
	}

	// Group by skill
	posList := make([][]int, S)
	Lvals := make([][]int64, S)
	Rvals := make([][]int64, S)
	skillID := make([]int, n)
	idxInSkill := make([]int, n)
	for i := 0; i < n; i++ {
		sid := mSkill[skills[i]]
		skillID[i] = sid
		idxInSkill[i] = len(posList[sid])
		posList[sid] = append(posList[sid], i)
		Lvals[sid] = append(Lvals[sid], L[i])
		Rvals[sid] = append(Rvals[sid], R[i])
	}

	trees := make([]*SegTree, S)
	for s := 0; s < S; s++ {
		if len(Lvals[s]) > 0 {
			trees[s] = BuildSegTree(Lvals[s], Rvals[s])
		}
	}

	active := make([]bool, n)
	for i := range active {
		active[i] = true
	}

	var total int64
	for _, t := range trees {
		if t != nil {
			total += t.Total()
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var typ, x int
		fmt.Fscan(in, &typ, &x)
		idx := x - 1
		sid := skillID[idx]
		ind := idxInSkill[idx]
		tree := trees[sid]
		old := tree.Total()
		if typ == 1 {
			if active[idx] {
				active[idx] = false
				tree.Update(ind, false)
			}
		} else {
			if !active[idx] {
				active[idx] = true
				tree.Update(ind, true)
			}
		}
		total += tree.Total() - old
		fmt.Fprintln(out, total%MOD)
	}
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

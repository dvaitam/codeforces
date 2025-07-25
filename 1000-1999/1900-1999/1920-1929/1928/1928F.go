package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) Add(i int, v int64) {
	for i <= f.n {
		f.tree[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

// treap implementation for ordered set of ints

type node struct {
	key         int
	pr          uint32
	left, right *node
}

type Treap struct{ root *node }

func rotateRight(p *node) *node {
	q := p.left
	p.left = q.right
	q.right = p
	return q
}
func rotateLeft(p *node) *node {
	q := p.right
	p.right = q.left
	q.left = p
	return q
}

func insertNode(p *node, key int) *node {
	if p == nil {
		return &node{key: key, pr: rand.Uint32()}
	}
	if key < p.key {
		p.left = insertNode(p.left, key)
		if p.left.pr < p.pr {
			p = rotateRight(p)
		}
	} else if key > p.key {
		p.right = insertNode(p.right, key)
		if p.right.pr < p.pr {
			p = rotateLeft(p)
		}
	}
	return p
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr < b.pr {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func eraseNode(p *node, key int) *node {
	if p == nil {
		return nil
	}
	if key < p.key {
		p.left = eraseNode(p.left, key)
	} else if key > p.key {
		p.right = eraseNode(p.right, key)
	} else {
		p = merge(p.left, p.right)
	}
	return p
}

func (t *Treap) Insert(key int) { t.root = insertNode(t.root, key) }
func (t *Treap) Remove(key int) { t.root = eraseNode(t.root, key) }

func (t *Treap) Prev(key int) int {
	cur := t.root
	res := -1
	for cur != nil {
		if key <= cur.key {
			cur = cur.left
		} else {
			res = cur.key
			cur = cur.right
		}
	}
	return res
}

func (t *Treap) Next(key int) int {
	cur := t.root
	res := 1<<31 - 1
	for cur != nil {
		if key >= cur.key {
			cur = cur.right
		} else {
			res = cur.key
			cur = cur.left
		}
	}
	return res
}

var (
	n, m                               int
	diffA                              []int64
	diffB                              []int64
	setA, setB                         Treap
	bitACnt, bitASum, bitASq, bitACube *Fenwick
	bitBCnt, bitBSum, bitBSq, bitBCube *Fenwick
	totalCntA, totalCntB               int64
	totalSumA, totalSumB               int64
	ans                                int64
)

func updateBits(cnt, sum, sq, cube *Fenwick, length int, delta int64) {
	cnt.Add(length, delta)
	sum.Add(length, delta*int64(length))
	sq.Add(length, delta*int64(length)*int64(length))
	cube.Add(length, delta*int64(length)*int64(length)*int64(length))
}

func contributionA(r int) int64 {
	if r == 0 {
		return 0
	}
	r64 := int64(r)
	cntLT := bitBCnt.Sum(r - 1)
	sumCLT := bitBSum.Sum(r - 1)
	sumSqLT := bitBSq.Sum(r - 1)
	sumCubeLT := bitBCube.Sum(r - 1)
	cntGE := totalCntB - cntLT
	sumCGE := totalSumB - sumCLT
	part1 := r64 * (r64 + 1) / 2 * sumCGE
	part2 := (-r64*r64*r64 + r64) / 6 * cntGE
	part3 := r64 * (sumSqLT + sumCLT) / 2
	part4 := (-sumCubeLT + sumCLT) / 6
	return part1 + part2 + part3 + part4
}

func contributionB(c int) int64 {
	if c == 0 {
		return 0
	}
	c64 := int64(c)
	cntLT := bitACnt.Sum(c - 1)
	sumCLT := bitASum.Sum(c - 1)
	sumSqLT := bitASq.Sum(c - 1)
	sumCubeLT := bitACube.Sum(c - 1)
	cntGE := totalCntA - cntLT
	sumCGE := totalSumA - sumCLT
	part1 := c64 * (c64 + 1) / 2 * sumCGE
	part2 := (-c64*c64*c64 + c64) / 6 * cntGE
	part3 := c64 * (sumSqLT + sumCLT) / 2
	part4 := (-sumCubeLT + sumCLT) / 6
	return part1 + part2 + part3 + part4
}

func addSegmentA(length int) {
	if length <= 0 {
		return
	}
	totalCntA++
	totalSumA += int64(length)
	updateBits(bitACnt, bitASum, bitASq, bitACube, length, 1)
	ans += contributionA(length)
}

func removeSegmentA(length int) {
	if length <= 0 {
		return
	}
	totalCntA--
	totalSumA -= int64(length)
	updateBits(bitACnt, bitASum, bitASq, bitACube, length, -1)
	ans -= contributionA(length)
}

func addSegmentB(length int) {
	if length <= 0 {
		return
	}
	totalCntB++
	totalSumB += int64(length)
	updateBits(bitBCnt, bitBSum, bitBSq, bitBCube, length, 1)
	ans += contributionB(length)
}

func removeSegmentB(length int) {
	if length <= 0 {
		return
	}
	totalCntB--
	totalSumB -= int64(length)
	updateBits(bitBCnt, bitBSum, bitBSq, bitBCube, length, -1)
	ans -= contributionB(length)
}

func modifyEdgeA(i int, delta int64) {
	old := diffA[i]
	newV := old + delta
	if old == 0 && newV != 0 {
		l := setA.Prev(i)
		r := setA.Next(i)
		removeSegmentA(i - l)
		removeSegmentA(r - i)
		addSegmentA(r - l)
		setA.Remove(i)
	} else if old != 0 && newV == 0 {
		l := setA.Prev(i)
		r := setA.Next(i)
		removeSegmentA(r - l)
		addSegmentA(i - l)
		addSegmentA(r - i)
		setA.Insert(i)
	}
	diffA[i] = newV
}

func modifyEdgeB(i int, delta int64) {
	old := diffB[i]
	newV := old + delta
	if old == 0 && newV != 0 {
		l := setB.Prev(i)
		r := setB.Next(i)
		removeSegmentB(i - l)
		removeSegmentB(r - i)
		addSegmentB(r - l)
		setB.Remove(i)
	} else if old != 0 && newV == 0 {
		l := setB.Prev(i)
		r := setB.Next(i)
		removeSegmentB(r - l)
		addSegmentB(i - l)
		addSegmentB(r - i)
		setB.Insert(i)
	}
	diffB[i] = newV
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &n, &m, &q)
	arrA := make([]int64, n)
	arrB := make([]int64, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arrA[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &arrB[i])
	}

	diffA = make([]int64, n)
	diffB = make([]int64, m)
	bitACnt = NewFenwick(n)
	bitASum = NewFenwick(n)
	bitASq = NewFenwick(n)
	bitACube = NewFenwick(n)
	bitBCnt = NewFenwick(m)
	bitBSum = NewFenwick(m)
	bitBSq = NewFenwick(m)
	bitBCube = NewFenwick(m)

	// build setA and segments
	setA.Insert(0)
	prev := 0
	for i := 1; i < n; i++ {
		diffA[i] = arrA[i] - arrA[i-1]
		if diffA[i] == 0 {
			setA.Insert(i)
			addSegmentA(i - prev)
			prev = i
		}
	}
	setA.Insert(n)
	addSegmentA(n - prev)

	// build setB
	setB.Insert(0)
	prev = 0
	for i := 1; i < m; i++ {
		diffB[i] = arrB[i] - arrB[i-1]
		if diffB[i] == 0 {
			setB.Insert(i)
			addSegmentB(i - prev)
			prev = i
		}
	}
	setB.Insert(m)
	addSegmentB(m - prev)

	fmt.Fprintln(writer, ans)
	for ; q > 0; q-- {
		var t, l, r int
		var x int64
		fmt.Fscan(reader, &t, &l, &r, &x)
		if t == 1 {
			if l > 1 {
				modifyEdgeA(l-1, x)
			}
			if r < n {
				modifyEdgeA(r, -x)
			}
		} else {
			if l > 1 {
				modifyEdgeB(l-1, x)
			}
			if r < m {
				modifyEdgeB(r, -x)
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

type treapNode struct {
	key   int
	prio  int
	left  *treapNode
	right *treapNode
}

func split(root *treapNode, key int) (l, r *treapNode) {
	if root == nil {
		return nil, nil
	}
	if root.key < key {
		root.right, r = split(root.right, key)
		l = root
	} else {
		l, root.left = split(root.left, key)
		r = root
	}
	return
}

func merge(l, r *treapNode) *treapNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		return l
	}
	r.left = merge(l, r.left)
	return r
}

func insert(root *treapNode, node *treapNode) *treapNode {
	if root == nil {
		return node
	}
	if node.prio > root.prio {
		node.left, node.right = split(root, node.key)
		return node
	}
	if node.key < root.key {
		root.left = insert(root.left, node)
	} else {
		root.right = insert(root.right, node)
	}
	return root
}

func deleteNode(root *treapNode, key int) *treapNode {
	if root == nil {
		return nil
	}
	if root.key == key {
		return merge(root.left, root.right)
	}
	if key < root.key {
		root.left = deleteNode(root.left, key)
	} else {
		root.right = deleteNode(root.right, key)
	}
	return root
}

func predecessor(root *treapNode, key int) int {
	res := -1 << 60
	for root != nil {
		if root.key < key {
			if root.key > res {
				res = root.key
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *treapNode, key int) int {
	res := 1<<60 - 1
	for root != nil {
		if root.key > key {
			if root.key < res {
				res = root.key
			}
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

type set struct{ root *treapNode }

func (s *set) insert(key int)   { s.root = insert(s.root, &treapNode{key: key, prio: rand.Int()}) }
func (s *set) remove(key int)   { s.root = deleteNode(s.root, key) }
func (s *set) prev(key int) int { return predecessor(s.root, key) }
func (s *set) next(key int) int { return successor(s.root, key) }

type segment struct {
	l, r int
	v    int64
}

var (
	n        int
	segments map[int]*segment
	s        set
	seed     int64
	vmax     int64
)

func rnd() int64 {
	ret := seed
	seed = (seed*7 + 13) % 1000000007
	return ret
}

func splitSeg(pos int) {
	if pos > n {
		return
	}
	start := s.prev(pos + 1)
	if start == pos {
		return
	}
	seg := segments[start]
	if pos > seg.r || pos <= seg.l {
		return
	}
	s.remove(start)
	delete(segments, start)
	left := &segment{l: start, r: pos - 1, v: seg.v}
	right := &segment{l: pos, r: seg.r, v: seg.v}
	segments[left.l] = left
	segments[right.l] = right
	s.insert(left.l)
	s.insert(right.l)
}

func rangeAdd(l, r int, x int64) {
	splitSeg(l)
	splitSeg(r + 1)
	for cur := l; cur <= r; {
		seg := segments[cur]
		seg.v += x
		next := s.next(cur)
		if next == cur { // should not happen
			break
		}
		cur = next
	}
}

func rangeAssign(l, r int, x int64) {
	splitSeg(l)
	splitSeg(r + 1)
	cur := l
	for cur <= r {
		next := s.next(cur)
		s.remove(cur)
		delete(segments, cur)
		cur = next
	}
	segments[l] = &segment{l: l, r: r, v: x}
	s.insert(l)
}

type pair struct {
	v int64
	c int64
}

func kth(l, r int, k int64) int64 {
	splitSeg(l)
	splitSeg(r + 1)
	var arr []pair
	for cur := l; cur <= r; {
		seg := segments[cur]
		arr = append(arr, pair{seg.v, int64(seg.r - seg.l + 1)})
		next := s.next(cur)
		cur = next
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].v < arr[j].v })
	for _, p := range arr {
		if k <= p.c {
			return p.v
		}
		k -= p.c
	}
	return -1
}

func powmod(a, b, mod int64) int64 {
	if mod == 1 {
		return 0
	}
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func rangePowSum(l, r int, x, mod int64) int64 {
	splitSeg(l)
	splitSeg(r + 1)
	ans := int64(0)
	for cur := l; cur <= r; {
		seg := segments[cur]
		cnt := int64(seg.r - seg.l + 1)
		ans = (ans + cnt*powmod(seg.v%mod, x, mod)) % mod
		next := s.next(cur)
		cur = next
	}
	return ans
}

func main() {
	rand.Seed(1)
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	fmt.Fscan(reader, &n, &m, &seed, &vmax)

	segments = make(map[int]*segment)
	for i := 1; i <= n; i++ {
		val := (rnd()%vmax + 1)
		segments[i] = &segment{l: i, r: i, v: val}
		s.insert(i)
	}

	for i := 0; i < m; i++ {
		op := int(rnd()%4) + 1
		l := int(rnd()%int64(n)) + 1
		r := int(rnd()%int64(n)) + 1
		if l > r {
			l, r = r, l
		}
		var x, y int64
		if op == 3 {
			x = rnd()%int64(r-l+1) + 1
		} else {
			x = rnd()%vmax + 1
		}
		if op == 4 {
			y = rnd()%vmax + 1
		}
		switch op {
		case 1:
			rangeAdd(l, r, x)
		case 2:
			rangeAssign(l, r, x)
		case 3:
			fmt.Fprintln(writer, kth(l, r, x))
		case 4:
			fmt.Fprintln(writer, rangePowSum(l, r, x, y))
		}
	}
}

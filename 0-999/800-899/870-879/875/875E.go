package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type node struct {
	key   int
	prio  int
	left  *node
	right *node
}

func split(t *node, key int) (l, r *node) {
	if t == nil {
		return nil, nil
	}
	if t.key < key {
		var sr *node
		t.right, sr = split(t.right, key)
		return t, sr
	}
	var sl *node
	sl, t.left = split(t.left, key)
	return sl, t
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.prio < b.prio {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

type Set struct{ root *node }

func (s *Set) Insert(key int) {
	l, r := split(s.root, key)
	m, rr := split(r, key+1)
	if m == nil {
		m = &node{key: key, prio: rand.Int()}
	}
	s.root = merge(merge(l, m), rr)
}

func (s *Set) Erase(key int) {
	l, r := split(s.root, key)
	_, rr := split(r, key+1)
	s.root = merge(l, rr)
}

func (s *Set) Next(key int) int {
	ans := -1
	for t := s.root; t != nil; {
		if t.key >= key {
			ans = t.key
			t = t.left
		} else {
			t = t.right
		}
	}
	return ans
}

func (s *Set) Prev(key int) int {
	ans := -1
	for t := s.root; t != nil; {
		if t.key <= key {
			ans = t.key
			t = t.right
		} else {
			t = t.left
		}
	}
	return ans
}

func (s *Set) Empty() bool { return s.root == nil }

func (s *Set) removeLessThan(limit int) {
	for {
		v := s.Prev(limit - 1)
		if v == -1 {
			break
		}
		s.Erase(v)
	}
}

func (s *Set) removeGreaterThan(limit int) {
	for {
		v := s.Next(limit + 1)
		if v == -1 {
			break
		}
		s.Erase(v)
	}
}

var (
	n      int
	s1, s2 int
	x      []int
)

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func check(d int64) bool {
	if abs64(int64(s1)-int64(s2)) > d {
		return false
	}
	set := &Set{}
	if abs64(int64(x[0])-int64(s2)) <= d {
		set.Insert(s2)
	}
	if abs64(int64(x[0])-int64(s1)) <= d {
		set.Insert(s1)
	}
	if set.Empty() {
		return false
	}
	for i := 0; i < n-1; i++ {
		prevNonEmpty := !set.Empty()
		L := int(int64(x[i+1]) - d)
		R := int(int64(x[i+1]) + d)
		set.removeLessThan(L)
		set.removeGreaterThan(R)
		if prevNonEmpty && abs64(int64(x[i+1])-int64(x[i])) <= d {
			set.Insert(x[i])
		}
		if set.Empty() {
			return false
		}
	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &s1, &s2)
	x = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}
	lo, hi := int64(0), int64(1_000_000_000)
	for lo < hi {
		mid := (lo + hi) / 2
		if check(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Print(lo)
}

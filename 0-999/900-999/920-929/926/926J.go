package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type treap struct {
	key         int
	prio        int
	left, right *treap
}

func split(t *treap, key int) (l, r *treap) {
	if t == nil {
		return nil, nil
	}
	if t.key < key {
		var sr *treap
		t.right, sr = split(t.right, key)
		return t, sr
	}
	var sl *treap
	sl, t.left = split(t.left, key)
	return sl, t
}

func merge(a, b *treap) *treap {
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

type Set struct{ root *treap }

func (s *Set) Insert(key int) {
	l, r := split(s.root, key)
	m, rr := split(r, key+1)
	if m == nil {
		m = &treap{key: key, prio: rand.Int()}
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

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	mp := make(map[int]int)
	s := &Set{}
	count := 0

	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		newL, newR := l, r

		// merge with intervals to the left
		prev := s.Prev(newL)
		for prev != -1 && mp[prev] >= newL {
			if mp[prev] > newR {
				newR = mp[prev]
			}
			newL = prev
			s.Erase(prev)
			delete(mp, prev)
			count--
			prev = s.Prev(newL)
		}
		// merge with intervals to the right
		next := s.Next(newL)
		for next != -1 && next <= newR {
			if mp[next] > newR {
				newR = mp[next]
			}
			s.Erase(next)
			delete(mp, next)
			count--
			next = s.Next(newL)
		}
		s.Insert(newL)
		mp[newL] = newR
		count++

		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, count)
	}
	out.WriteByte('\n')
}

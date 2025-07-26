package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"sort"
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
	var l, r = split(s.root, key)
	var m, rr = split(r, key+1)
	if m == nil {
		m = &treap{key: key, prio: rand.Int()}
	}
	s.root = merge(merge(l, m), rr)
}

func (s *Set) Erase(key int) {
	var l, r = split(s.root, key)
	var m, rr = split(r, key+1)
	_ = m
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

func (s *Set) Min() (int, bool) {
	t := s.root
	if t == nil {
		return 0, false
	}
	for t.left != nil {
		t = t.left
	}
	return t.key, true
}

func (s *Set) Max() (int, bool) {
	t := s.root
	if t == nil {
		return 0, false
	}
	for t.right != nil {
		t = t.right
	}
	return t.key, true
}

// max heap for gaps
type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h IntHeap) Peek() int {
	if len(h) == 0 {
		return 0
	}
	return h[0]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var q int
		fmt.Fscan(reader, &q)

		freq := make(map[int]int)
		set := &Set{}
		var gaps IntHeap
		gapCnt := make(map[int]int)
		heap.Init(&gaps)

		// insert initial values
		for _, v := range a {
			freq[v]++
		}
		vals := make([]int, 0, len(freq))
		for v := range freq {
			vals = append(vals, v)
		}
		sort.Ints(vals)
		for _, v := range vals {
			set.Insert(v)
		}
		for i := 0; i+1 < len(vals); i++ {
			g := vals[i+1] - vals[i]
			gapCnt[g]++
			heap.Push(&gaps, g)
		}
		uniqueCount := len(vals)

		getMaxGap := func() int {
			for gaps.Len() > 0 {
				g := gaps.Peek()
				if gapCnt[g] > 0 {
					return g
				}
				heap.Pop(&gaps)
			}
			return 0
		}

		for ; q > 0; q-- {
			var idx, x int
			fmt.Fscan(reader, &idx, &x)
			idx--
			old := a[idx]
			if old != x {
				// remove old value
				if freq[old] == 1 {
					pred := set.Prev(old - 1)
					succ := set.Next(old + 1)
					if pred != -1 {
						g := old - pred
						gapCnt[g]--
						heap.Push(&gaps, g)
					}
					if succ != -1 {
						g := succ - old
						gapCnt[g]--
						heap.Push(&gaps, g)
					}
					if pred != -1 && succ != -1 {
						g := succ - pred
						gapCnt[g]++
						heap.Push(&gaps, g)
					}
					set.Erase(old)
					uniqueCount--
				}
				freq[old]--
				if freq[old] == 0 {
					delete(freq, old)
				}

				// insert new value
				if freq[x] == 0 {
					pred := set.Prev(x - 1)
					succ := set.Next(x + 1)
					if pred != -1 && succ != -1 {
						g := succ - pred
						gapCnt[g]--
						heap.Push(&gaps, g)
					}
					if pred != -1 {
						g := x - pred
						gapCnt[g]++
						heap.Push(&gaps, g)
					}
					if succ != -1 {
						g := succ - x
						gapCnt[g]++
						heap.Push(&gaps, g)
					}
					set.Insert(x)
					uniqueCount++
				}
				freq[x]++
				a[idx] = x
			}
			if uniqueCount == 1 {
				val, _ := set.Max()
				fmt.Fprintln(writer, val)
			} else {
				maxVal, _ := set.Max()
				maxGap := getMaxGap()
				fmt.Fprintln(writer, maxVal+maxGap)
			}
		}
	}
}

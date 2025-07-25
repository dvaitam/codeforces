package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type SegMax struct {
	n   int
	arr []int
}

func NewSegMax(n int) *SegMax {
	m := 1
	for m < n {
		m <<= 1
	}
	return &SegMax{n: m, arr: make([]int, 2*m)}
}

func (s *SegMax) Update(pos, val int) {
	pos += s.n
	if val > s.arr[pos] {
		s.arr[pos] = val
	}
	for pos >>= 1; pos > 0; pos >>= 1 {
		left := s.arr[pos<<1]
		right := s.arr[pos<<1|1]
		if left > right {
			s.arr[pos] = left
		} else {
			s.arr[pos] = right
		}
	}
}

func (s *SegMax) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += s.n
	r += s.n
	res := 0
	for l <= r {
		if l&1 == 1 {
			if s.arr[l] > res {
				res = s.arr[l]
			}
			l++
		}
		if r&1 == 0 {
			if s.arr[r] > res {
				res = s.arr[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

type SegMin struct {
	n   int
	inf int
	arr []int
}

func NewSegMin(n int) *SegMin {
	m := 1
	for m < n {
		m <<= 1
	}
	inf := int(1e9)
	arr := make([]int, 2*m)
	for i := range arr {
		arr[i] = inf
	}
	return &SegMin{n: m, inf: inf, arr: arr}
}

func (s *SegMin) Update(pos, val int) {
	pos += s.n
	if val < s.arr[pos] {
		s.arr[pos] = val
	}
	for pos >>= 1; pos > 0; pos >>= 1 {
		left := s.arr[pos<<1]
		right := s.arr[pos<<1|1]
		if left < right {
			s.arr[pos] = left
		} else {
			s.arr[pos] = right
		}
	}
}

func (s *SegMin) Query(l, r int) int {
	if l > r {
		return s.inf
	}
	l += s.n
	r += s.n
	res := s.inf
	for l <= r {
		if l&1 == 1 {
			if s.arr[l] < res {
				res = s.arr[l]
			}
			l++
		}
		if r&1 == 0 {
			if s.arr[r] < res {
				res = s.arr[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		ok := true
		maxVal := n
		for i := 0; i < n; i++ {
			if b[i] < a[i] {
				ok = false
			}
			if a[i] > maxVal {
				maxVal = a[i]
			}
			if b[i] > maxVal {
				maxVal = b[i]
			}
		}
		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}

		segA := NewSegMax(maxVal + 2)
		segB := NewSegMax(maxVal + 2)
		left := make([]int, n)
		for i := 0; i < n; i++ {
			x := b[i]
			l1 := segA.Query(x+1, maxVal+1)
			l2 := segB.Query(0, x-1)
			if l1 > l2 {
				left[i] = l1
			} else {
				left[i] = l2
			}
			segA.Update(a[i], i+1)
			segB.Update(b[i], i+1)
		}

		segAR := NewSegMin(maxVal + 2)
		segBR := NewSegMin(maxVal + 2)
		right := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			x := b[i]
			r1 := segAR.Query(x+1, maxVal+1)
			r2 := segBR.Query(0, x-1)
			if r1 < r2 {
				right[i] = r1
			} else {
				right[i] = r2
			}
			segAR.Update(a[i], i+1)
			segBR.Update(b[i], i+1)
		}
		for i := 0; i < n; i++ {
			if right[i] == segAR.inf {
				right[i] = n + 1
			}
		}

		pos := make(map[int][]int)
		for i, v := range a {
			pos[v] = append(pos[v], i+1)
		}

		for i, x := range b {
			L := left[i] + 1
			R := right[i] - 1
			arr := pos[x]
			idx := sort.SearchInts(arr, L)
			if idx == len(arr) || arr[idx] > R {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

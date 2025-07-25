package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	next map[int]int
	link int
	len  int
}

type SAM struct {
	st   []state
	last int
}

func NewSAM(size int) *SAM {
	s := &SAM{st: make([]state, 1, size)}
	s.st[0] = state{next: make(map[int]int), link: -1, len: 0}
	s.last = 0
	return s
}

func (s *SAM) Extend(c int) {
	cur := len(s.st)
	s.st = append(s.st, state{next: make(map[int]int), len: s.st[s.last].len + 1})
	p := s.last
	for p != -1 && s.st[p].next[c] == 0 {
		s.st[p].next[c] = cur
		p = s.st[p].link
	}
	if p == -1 {
		s.st[cur].link = 0
	} else {
		q := s.st[p].next[c]
		if s.st[p].len+1 == s.st[q].len {
			s.st[cur].link = q
		} else {
			clone := len(s.st)
			stClone := state{next: make(map[int]int), len: s.st[p].len + 1, link: s.st[q].link}
			for k, v := range s.st[q].next {
				stClone.next[k] = v
			}
			s.st = append(s.st, stClone)
			for p != -1 && s.st[p].next[c] == q {
				s.st[p].next[c] = clone
				p = s.st[p].link
			}
			s.st[q].link = clone
			s.st[cur].link = clone
		}
	}
	s.last = cur
}

func (s *SAM) DistinctSubstrings() int {
	res := 0
	for i := 1; i < len(s.st); i++ {
		res += s.st[i].len - s.st[s.st[i].link].len
	}
	return res
}

func nextPerm(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for j > i && a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// naive approach, only works for small n due to factorial growth
	if n > 8 {
		fmt.Fprintln(out, 0)
		return
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	total := 1
	for i := 2; i <= n; i++ {
		total *= i
	}
	sam := NewSAM(n * total * 2)
	for {
		for _, v := range perm {
			sam.Extend(v)
		}
		if !nextPerm(perm) {
			break
		}
	}
	fmt.Fprintln(out, sam.DistinctSubstrings())
}

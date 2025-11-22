package main

import (
	"bufio"
	"fmt"
	"os"
)

type solver struct {
	n      int
	a      []int
	target []int
	mis    [3][]int // indices of positions that are incorrect and currently have value v
	inMis  []bool
	bad    int // total incorrect positions
}

func (s *solver) boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// update refreshes membership of index i in mis lists after a value change.
func (s *solver) update(i int) {
	if s.a[i] == s.target[i] {
		s.inMis[i] = false
		return
	}
	if !s.inMis[i] {
		s.inMis[i] = true
		s.mis[s.a[i]] = append(s.mis[s.a[i]], i)
		return
	}
	// Already marked incorrect but value could have changed; keep the flag
	// and add the index to the new value's bucket.
	s.mis[s.a[i]] = append(s.mis[s.a[i]], i)
}

func (s *solver) popValue(val, avoid int) int {
	for len(s.mis[val]) > 0 {
		last := len(s.mis[val]) - 1
		idx := s.mis[val][last]
		s.mis[val] = s.mis[val][:last]
		if idx == avoid {
			continue
		}
		if s.inMis[idx] && s.a[idx] == val && s.a[idx] != s.target[idx] {
			return idx
		}
	}
	// The problem guarantees a solution, so we should never get here.
	panic("no index with required value")
}

func (s *solver) popAny(avoid int) int {
	for _, v := range []int{0, 2} { // only need indices with value differing from 1
		for len(s.mis[v]) > 0 {
			last := len(s.mis[v]) - 1
			idx := s.mis[v][last]
			s.mis[v] = s.mis[v][:last]
			if idx == avoid {
				continue
			}
			if s.inMis[idx] && s.a[idx] != s.target[idx] {
				return idx
			}
		}
	}
	panic("no incorrect index found")
}

func (s *solver) solve() []struct{ u, v int } {
	cnt := [3]int{}
	for _, x := range s.a {
		cnt[x]++
	}

	s.target = make([]int, s.n)
	for i := 0; i < cnt[0]; i++ {
		s.target[i] = 0
	}
	for i := cnt[0]; i < cnt[0]+cnt[1]; i++ {
		s.target[i] = 1
	}
	for i := cnt[0] + cnt[1]; i < s.n; i++ {
		s.target[i] = 2
	}

	s.inMis = make([]bool, s.n)
	s.bad = 0
	for i := 0; i < s.n; i++ {
		if s.a[i] != s.target[i] {
			s.inMis[i] = true
			s.bad++
			s.mis[s.a[i]] = append(s.mis[s.a[i]], i)
		}
	}

	if s.bad == 0 {
		return nil
	}

	ops := make([]struct{ u, v int }, 0)

	pos := -1
	for i, x := range s.a {
		if x == 1 {
			pos = i
			break
		}
	}
	if pos == -1 {
		panic("no value 1 found, contradicts constraints")
	}

	for s.bad > 0 {
		var idx int
		if s.target[pos] == 1 {
			idx = s.popAny(pos)
		} else {
			idx = s.popValue(s.target[pos], pos)
		}

		beforePos := s.a[pos] != s.target[pos]
		beforeIdx := s.a[idx] != s.target[idx]

		// Remove stale membership before values change.
		s.inMis[pos] = false
		s.inMis[idx] = false

		s.a[pos], s.a[idx] = s.a[idx], s.a[pos]

		afterPos := s.a[pos] != s.target[pos]
		afterIdx := s.a[idx] != s.target[idx]
		s.bad += s.boolToInt(afterPos) + s.boolToInt(afterIdx) - s.boolToInt(beforePos) - s.boolToInt(beforeIdx)

		s.update(pos)
		s.update(idx)

		ops = append(ops, struct{ u, v int }{u: pos + 1, v: idx + 1})
		pos = idx
	}

	return ops
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		s := solver{n: n, a: a}
		ops := s.solve()

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintf(out, "%d %d\n", op.u, op.v)
		}
	}
}

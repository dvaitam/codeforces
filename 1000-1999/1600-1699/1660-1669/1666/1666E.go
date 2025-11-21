package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type solver struct {
	l int64
	n int
	a []int64
}

func newSolver(l int64, a []int64) *solver {
	return &solver{
		l: l,
		n: len(a),
		a: a,
	}
}

func (s *solver) propagate(L, D int64) (bool, int64, int64, bool) {
	if s.n == 1 {
		return true, 0, 0, false
	}
	lo := int64(0)
	hi := int64(0)
	for i := 0; i < s.n-1; i++ {
		left := s.a[i]
		right := s.a[i+1]
		tmpLo := lo + L
		if tmpLo < left {
			tmpLo = left
		}
		lo = tmpLo
		tmpHiCandidate := hi + L + D
		tmpHi := right
		if tmpHiCandidate < tmpHi {
			tmpHi = tmpHiCandidate
		}
		hi = tmpHi
		if lo > hi {
			needIncrease := tmpHiCandidate < left
			return false, lo, hi, needIncrease
		}
	}
	return true, lo, hi, false
}

func (s *solver) propagateWithStore(L, D int64) (bool, []int64, []int64) {
	loArr := make([]int64, s.n)
	hiArr := make([]int64, s.n)
	lo := int64(0)
	hi := int64(0)
	loArr[0] = 0
	hiArr[0] = 0
	if s.n == 1 {
		return true, loArr, hiArr
	}
	for i := 1; i <= s.n-1; i++ {
		left := s.a[i-1]
		right := s.a[i]
		tmpLo := lo + L
		if tmpLo < left {
			tmpLo = left
		}
		lo = tmpLo
		tmpHiCandidate := hi + L + D
		tmpHi := right
		if tmpHiCandidate < tmpHi {
			tmpHi = tmpHiCandidate
		}
		hi = tmpHi
		loArr[i] = lo
		hiArr[i] = hi
		if lo > hi {
			return false, nil, nil
		}
	}
	return true, loArr, hiArr
}

func (s *solver) findFirstOk(Lmax, D int64) int64 {
	lo := int64(1)
	hi := Lmax
	ans := int64(-1)
	for lo <= hi {
		mid := (lo + hi) >> 1
		ok, _, _, needIncrease := s.propagate(mid, D)
		if ok {
			ans = mid
			hi = mid - 1
		} else {
			if needIncrease {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return ans
}

func (s *solver) feasible(D int64, needSegments bool) (bool, []int64) {
	if s.n == 1 {
		if needSegments {
			return true, []int64{0, s.l}
		}
		return true, nil
	}
	Lmax := s.l / int64(s.n)
	if Lmax == 0 {
		return false, nil
	}
	first := s.findFirstOk(Lmax, D)
	if first == -1 {
		return false, nil
	}
	bestL := int64(-1)
	bestLo := int64(0)
	lo := first
	hi := Lmax
	for lo <= hi {
		mid := (lo + hi) >> 1
		ok, curLo, curHi, needIncrease := s.propagate(mid, D)
		if !ok {
			if needIncrease {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
			continue
		}
		lowBound := s.l - (mid + D)
		if curHi >= lowBound {
			bestL = mid
			bestLo = curLo
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	if bestL == -1 {
		return false, nil
	}
	highBound := s.l - bestL
	if bestLo > highBound {
		return false, nil
	}
	if !needSegments {
		return true, nil
	}
	ok, loArr, hiArr := s.propagateWithStore(bestL, D)
	if !ok {
		return false, nil
	}
	lowTarget := max(loArr[s.n-1], s.l-(bestL+D))
	highTarget := min(hiArr[s.n-1], s.l-bestL)
	target := s.l - bestL
	if target < lowTarget {
		target = lowTarget
	} else if target > highTarget {
		target = highTarget
	}
	positions := make([]int64, s.n+1)
	positions[0] = 0
	positions[s.n] = s.l
	positions[s.n-1] = target
	for k := s.n - 1; k >= 1; k-- {
		minAllowed := positions[k] - (bestL + D)
		maxAllowed := positions[k] - bestL
		if minAllowed < loArr[k-1] {
			minAllowed = loArr[k-1]
		}
		if maxAllowed > hiArr[k-1] {
			maxAllowed = hiArr[k-1]
		}
		if maxAllowed < minAllowed {
			maxAllowed = minAllowed
		}
		positions[k-1] = maxAllowed
	}
	return true, positions
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var l int64
	var n int
	if _, err := fmt.Fscan(in, &l, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	s := newSolver(l, a)
	low := int64(0)
	high := l
	ansD := l
	for low <= high {
		mid := (low + high) >> 1
		ok, _ := s.feasible(mid, false)
		if ok {
			ansD = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	ok, positions := s.feasible(ansD, true)
	if !ok {
		return
	}
	for i := 1; i <= n; i++ {
		fmt.Fprintf(out, "%d %d\n", positions[i-1], positions[i])
	}
}

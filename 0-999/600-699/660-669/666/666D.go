package main

import (
	"bufio"
	"fmt"
	"os"
)

const BIG int64 = 2000000000000000000
const INF int64 = 8000000000000000000

type Point struct {
	x, y int64
}

type Side struct {
	fixed bool
	f     int64
	has   bool
	l, r  int64
}

type Sol struct {
	u, d, s int64
	lab     [4]int
}

var perms [][4]int
var xside = [4]int{0, 1, 0, 1}
var yside = [4]int{0, 0, 1, 1}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func gen(idx int, used [4]bool, cur [4]int) {
	if idx == 4 {
		perms = append(perms, cur)
		return
	}
	for i := 0; i < 4; i++ {
		if !used[i] {
			used[i] = true
			cur[idx] = i
			gen(idx+1, used, cur)
			used[i] = false
		}
	}
}

func addFix(s *Side, v int64) bool {
	if s.fixed && s.f != v {
		return false
	}
	s.fixed = true
	s.f = v
	return true
}

func addInter(s *Side, l, r int64) bool {
	if !s.has {
		s.has = true
		s.l = l
		s.r = r
	} else {
		if l > s.l {
			s.l = l
		}
		if r < s.r {
			s.r = r
		}
	}
	return s.l <= s.r
}

func sRange(a, b Side) (bool, int64, int64) {
	if (a.has && a.l > a.r) || (b.has && b.l > b.r) {
		return false, 0, 0
	}
	if a.fixed {
		if a.has && (a.f < a.l || a.f > a.r) {
			return false, 0, 0
		}
		if b.fixed {
			if b.has && (b.f < b.l || b.f > b.r) {
				return false, 0, 0
			}
			s := b.f - a.f
			return true, s, s
		}
		if b.has {
			return true, b.l - a.f, b.r - a.f
		}
		return true, -INF, INF
	}
	if b.fixed {
		if b.has && (b.f < b.l || b.f > b.r) {
			return false, 0, 0
		}
		if a.has {
			return true, b.f - a.r, b.f - a.l
		}
		return true, -INF, INF
	}
	if a.has && b.has {
		return true, b.l - a.r, b.r - a.l
	}
	return true, -INF, INF
}

func chooseSide(s int64, a, b Side) (bool, int64) {
	lo, hi := int64(-INF), int64(INF)
	if a.has {
		if a.l > lo {
			lo = a.l
		}
		if a.r < hi {
			hi = a.r
		}
	}
	if a.fixed {
		if a.f > lo {
			lo = a.f
		}
		if a.f < hi {
			hi = a.f
		}
	}
	if b.has {
		l := b.l - s
		r := b.r - s
		if l > lo {
			lo = l
		}
		if r < hi {
			hi = r
		}
	}
	if b.fixed {
		v := b.f - s
		if v > lo {
			lo = v
		}
		if v < hi {
			hi = v
		}
	}
	if lo > hi {
		return false, 0
	}
	if lo <= 0 && 0 <= hi {
		return true, 0
	}
	if lo > 0 {
		return true, lo
	}
	return true, hi
}

func caseFeasible(p [4]Point, perm [4]int, mask int, L int64) (bool, int64, int64, int64) {
	var u, v, d, e Side
	for i := 0; i < 4; i++ {
		j := perm[i]
		if ((mask >> i) & 1) == 1 {
			if xside[j] == 0 {
				if !addFix(&u, p[i].x) {
					return false, 0, 0, 0
				}
			} else {
				if !addFix(&v, p[i].x) {
					return false, 0, 0, 0
				}
			}
			l := p[i].y - L
			r := p[i].y + L
			if yside[j] == 0 {
				if !addInter(&d, l, r) {
					return false, 0, 0, 0
				}
			} else {
				if !addInter(&e, l, r) {
					return false, 0, 0, 0
				}
			}
		} else {
			if yside[j] == 0 {
				if !addFix(&d, p[i].y) {
					return false, 0, 0, 0
				}
			} else {
				if !addFix(&e, p[i].y) {
					return false, 0, 0, 0
				}
			}
			l := p[i].x - L
			r := p[i].x + L
			if xside[j] == 0 {
				if !addInter(&u, l, r) {
					return false, 0, 0, 0
				}
			} else {
				if !addInter(&v, l, r) {
					return false, 0, 0, 0
				}
			}
		}
	}

	okx, xl, xr := sRange(u, v)
	if !okx {
		return false, 0, 0, 0
	}
	oky, yl, yr := sRange(d, e)
	if !oky {
		return false, 0, 0, 0
	}

	lo := max64(max64(xl, yl), 1)
	hi := min64(xr, yr)
	if lo > hi {
		return false, 0, 0, 0
	}
	s := lo

	oku, uu := chooseSide(s, u, v)
	if !oku {
		return false, 0, 0, 0
	}
	okd, dd := chooseSide(s, d, e)
	if !okd {
		return false, 0, 0, 0
	}

	return true, uu, dd, s
}

func feasible(p [4]Point, L int64) (bool, Sol) {
	for _, perm := range perms {
		for mask := 0; mask < 16; mask++ {
			ok, u, d, s := caseFeasible(p, perm, mask, L)
			if ok {
				return true, Sol{u, d, s, perm}
			}
		}
	}
	return false, Sol{}
}

func main() {
	gen(0, [4]bool{}, [4]int{})

	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var p [4]Point
		for i := 0; i < 4; i++ {
			fmt.Fscan(in, &p[i].x, &p[i].y)
		}

		ok, _ := feasible(p, BIG)
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		lo, hi := int64(-1), BIG
		for hi-lo > 1 {
			mid := lo + (hi-lo)/2
			ok, _ := feasible(p, mid)
			if ok {
				hi = mid
			} else {
				lo = mid
			}
		}

		_, sol := feasible(p, hi)
		fmt.Fprintln(out, hi)
		for i := 0; i < 4; i++ {
			j := sol.lab[i]
			x := sol.u
			y := sol.d
			if xside[j] == 1 {
				x += sol.s
			}
			if yside[j] == 1 {
				y += sol.s
			}
			fmt.Fprintln(out, x, y)
		}
	}
}

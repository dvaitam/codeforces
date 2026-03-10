package main

import (
	"bufio"
	"fmt"
	"os"
)

var ps []int64

func getCost(l, r int, paid int64) int64 {
	for r > l+1 && ps[l]+paid > ps[l+1] {
		l++
	}
	for r > l+1 && ps[r]-paid < ps[r-1] {
		r--
	}
	if r <= l+1 {
		v := ps[r] - ps[l] + 1 - paid
		if v <= 0 {
			return 0
		}
		return v
	}

	// x0: argmin over i in [l+1,r] of slope from l: (ps[i]-ps[l])/(i-l)
	x0 := l + 1
	for i := l + 2; i <= r; i++ {
		// slope(i) < slope(x0) iff (ps[i]-ps[l])*(x0-l) < (ps[x0]-ps[l])*(i-l)
		if (ps[i]-ps[l])*int64(x0-l) < (ps[x0]-ps[l])*int64(i-l) {
			x0 = i
		}
	}
	if ps[x0] < ps[l]+paid*int64(x0-l) {
		return getCost(l, x0, paid) + getCost(x0, r, paid)
	}

	// x1: argmin over i in [l,r-1] of slope to r: (ps[r]-ps[i])/(r-i)
	x1 := r - 1
	for i := r - 2; i >= l; i-- {
		// slope_r(i) < slope_r(x1) iff (ps[r]-ps[i])*(r-x1) < (ps[r]-ps[x1])*(r-i)
		if (ps[r]-ps[i])*int64(r-x1) < (ps[r]-ps[x1])*int64(r-i) {
			x1 = i
		}
	}
	if ps[x1] > ps[r]-paid*int64(r-x1) {
		return getCost(l, x1, paid) + getCost(x1, r, paid)
	}

	s0 := (ps[x0] - ps[l]) / int64(x0-l)
	s1 := (ps[r] - ps[x1]) / int64(r-x1)
	cur := s0
	if s1 < cur {
		cur = s1
	}
	cur = cur + 1 - paid
	return cur + getCost(l, r, paid+cur)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	ps = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &ps[i])
		ps[i] += ps[i-1]
	}
	fmt.Println(getCost(0, n, 0))
}

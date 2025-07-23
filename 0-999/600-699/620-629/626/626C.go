package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int64
	fmt.Fscan(in, &n, &m)
	l, r := int64(0), int64(6*(n+m))
	if r == 0 {
		r = 1
	}
	feasible := func(x int64) bool {
		twosOnly := x/2 - x/6
		threesOnly := x/3 - x/6
		six := x / 6
		need2 := n - twosOnly
		if need2 < 0 {
			need2 = 0
		}
		need3 := m - threesOnly
		if need3 < 0 {
			need3 = 0
		}
		return need2+need3 <= six
	}
	for l+1 < r {
		mid := (l + r) / 2
		if feasible(mid) {
			r = mid
		} else {
			l = mid
		}
	}
	fmt.Println(r)
}

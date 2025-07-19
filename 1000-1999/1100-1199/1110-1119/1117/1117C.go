package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x1, y1, x2, y2 int64
	var n int
	if _, err := fmt.Fscan(reader, &x1, &y1, &x2, &y2, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	// relative target
	x2 -= x1
	y2 -= y1
	// net change per full cycle on target vector
	var dx, dy int64
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'U':
			dy--
		case 'D':
			dy++
		case 'L':
			dx++
		case 'R':
			dx--
		}
	}
	// binary search on full cycles k: find max k where still unreachable
	var l, r int64 = 1, 2000000000
	for l <= r {
		mid := (l + r) / 2
		tx := x2 + dx*mid
		ty := y2 + dy*mid
		// distance after mid cycles minus steps available
		tmp := abs64(tx) + abs64(ty) - mid*int64(n)
		if tmp <= 0 {
			// reachable in <= mid cycles
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	// if too many cycles still unreachable
	const thr = int64(1000000000 + 100000000)
	if r >= thr {
		fmt.Println(-1)
		return
	}
	// simulate r full cycles
	total := r * int64(n)
	x2 += dx * r
	y2 += dy * r
	// try partial cycle
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'U':
			y2--
		case 'D':
			y2++
		case 'L':
			x2++
		case 'R':
			x2--
		}
		if abs64(x2)+abs64(y2) <= total+int64(i+1) {
			fmt.Println(total + int64(i+1))
			return
		}
	}
	// should not reach here
	fmt.Println(-1)
}

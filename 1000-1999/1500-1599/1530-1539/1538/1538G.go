package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func feasible(k, x, y, a, b int64) bool {
	if a < b {
		a, b = b, a
		x, y = y, x
	}
	if a == b {
		return k <= min(x, y)/a
	}
	diff := a - b
	if k*(a+b) > x+y {
		return false
	}
	lower := (k*a - y + diff - 1) / diff
	if lower < 0 {
		lower = 0
	}
	upper := (x - k*b) / diff
	if upper > k {
		upper = k
	}
	if upper < 0 {
		return false
	}
	return lower <= upper
}
func maxSets(x, y, a, b int64) int64 {
	if x < y {
		x, y = y, x
	}
	if a < b {
		a, b = b, a
	}
	if a == b {
		return min(x, y) / a
	}
	hi := (x + y) / (a + b)
	lo := int64(0)
	var res int64
	for lo <= hi {
		mid := (lo + hi) / 2
		if feasible(mid, x, y, a, b) {
			res = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return res
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, a, b int64
		fmt.Fscan(in, &x, &y, &a, &b)
		fmt.Fprintln(out, maxSets(x, y, a, b))
	}
}

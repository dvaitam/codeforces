package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(x, length int64) int64 {
	if length <= x-1 {
		return (x - 1 + x - length) * length / 2
	}
	// x-1 positions decreasing to 1, rest ones
	return (x-1)*x/2 + (length - (x - 1))
}

func feasible(n, m, k, x int64) bool {
	left := calc(x, k-1)
	right := calc(x, n-k)
	return x+left+right <= m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int64
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	l, r := int64(1), m
	var ans int64
	for l <= r {
		mid := (l + r) / 2
		if feasible(n, m, k, mid) {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	fmt.Println(ans)
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func canComplete(segs [][2]int, k int) bool {
	low, high := 0, 0
	for _, seg := range segs {
		l, r := seg[0], seg[1]
		if low-k > l {
			l = low - k
		}
		if high+k < r {
			r = high + k
		}
		if l > r {
			return false
		}
		low, high = l, r
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		segs := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &segs[i][0], &segs[i][1])
		}

		lo, hi := 0, int(1e9)
		for lo < hi {
			mid := (lo + hi) / 2
			if canComplete(segs, mid) {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(writer, lo)
	}
}

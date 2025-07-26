package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(out *bufio.Writer, in *bufio.Reader, w int) int {
	fmt.Fprintf(out, "? %d\n", w)
	out.Flush()
	var h int
	fmt.Fscan(in, &h)
	return h
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	// binary search minimal width that results in height 1
	left, right := 1, 2000*(n+1)
	ans := right
	for left <= right {
		mid := (left + right) / 2
		h := query(out, in, mid)
		if h == 0 {
			left = mid + 1
		} else {
			ans = min(ans, h*mid)
			if h == 1 {
				right = mid - 1
			} else {
				right = mid - 1
			}
		}
	}

	// try other possible heights
	for h := 2; h <= n; h++ {
		w := ans / h
		if w == 0 {
			break
		}
		res := query(out, in, w)
		if res > 0 {
			ans = min(ans, res*w)
		}
	}

	fmt.Fprintf(out, "! %d\n", ans)
	out.Flush()
}

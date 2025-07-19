package main

import (
	"fmt"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	var W, H int64
	if _, err := fmt.Scan(&W, &H); err != nil {
		return
	}
	var x1, y1, x2, y2 int64
	fmt.Scan(&x1, &y1, &x2, &y2)
	var w, h int64
	fmt.Scan(&w, &h)
	const inf = int64(1e18)
	ans := inf
	// Try placing horizontally
	if (x2 - x1 + w) <= W {
		// place new table to left of old
		ans = min(ans, max(0, w-x1))
		// place new table to right of old
		ans = min(ans, max(0, x2-(W-w)))
	}
	// Try placing vertically
	if (y2 - y1 + h) <= H {
		// place new table below old
		ans = min(ans, max(0, h-y1))
		// place new table above old
		ans = min(ans, max(0, y2-(H-h)))
	}
	if ans == inf {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}

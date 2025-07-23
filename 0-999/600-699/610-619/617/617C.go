package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var x1, y1, x2, y2 int
	if _, err := fmt.Fscan(in, &n, &x1, &y1, &x2, &y2); err != nil {
		return
	}
	d1 := make([]int64, n)
	d2 := make([]int64, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		dx1 := int64(x - x1)
		dy1 := int64(y - y1)
		d1[i] = dx1*dx1 + dy1*dy1
		dx2 := int64(x - x2)
		dy2 := int64(y - y2)
		d2[i] = dx2*dx2 + dy2*dy2
	}
	const inf int64 = 1 << 62
	ans := inf
	// consider r1 determined by each flower
	for i := 0; i < n; i++ {
		r1 := d1[i]
		var r2 int64
		for j := 0; j < n; j++ {
			if d1[j] > r1 {
				if d2[j] > r2 {
					r2 = d2[j]
				}
			}
		}
		if r1+r2 < ans {
			ans = r1 + r2
		}
	}
	// also consider r1 = 0 (no radius at first fountain)
	var r2 int64
	for i := 0; i < n; i++ {
		if d2[i] > r2 {
			r2 = d2[i]
		}
	}
	if r2 < ans {
		ans = r2
	}
	fmt.Println(ans)
}

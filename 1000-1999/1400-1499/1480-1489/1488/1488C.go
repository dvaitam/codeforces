package main

import (
	"bufio"
	"fmt"
	"os"
)

func timeToCover(s, l, r int) int {
	if l > r {
		return 0
	}
	if l < 1 {
		l = 1
	}
	if r < 1 {
		r = 1
	}
	// r>=l
	interval := r - l
	d1 := s - l
	if d1 < 0 {
		d1 = -d1
	}
	d2 := s - r
	if d2 < 0 {
		d2 = -d2
	}
	if d1 < d2 {
		return interval + d1
	}
	return interval + d2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, x, y int
		fmt.Fscan(in, &n, &x, &y)
		if x > y {
			x, y = y, x
		}
		best := int(^uint(0) >> 1) // max int
		for m := 0; m <= n; m++ {
			t1 := timeToCover(x, 1, m)
			t2 := timeToCover(y, m+1, n)
			if t1 > t2 {
				if t1 < best {
					best = t1
				}
			} else {
				if t2 < best {
					best = t2
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}

package main

import (
	"fmt"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	var l, r, x, y int64
	if _, err := fmt.Scan(&l, &r, &x, &y); err != nil {
		return
	}
	if y%x != 0 {
		fmt.Println(0)
		return
	}
	k := y / x
	var cnt int
	for d := int64(1); d*d <= k; d++ {
		if k%d == 0 {
			m := d
			n := k / d
			if gcd(m, n) == 1 {
				a := x * m
				b := x * n
				if a >= l && a <= r && b >= l && b <= r {
					if m == n {
						cnt++
					} else {
						cnt += 2
					}
				}
			}
		}
	}
	fmt.Println(cnt)
}

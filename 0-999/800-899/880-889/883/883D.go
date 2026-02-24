package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scan(&n)
	var s string
	fmt.Scan(&s)

	pre := make([]int, n+2)
	d := make([]int, 1)
	p := make([]int, 1)

	for i := 1; i <= n; i++ {
		ch := s[i-1]
		if ch == '*' {
			pre[i] = pre[i-1] + 1
			d = append(d, i)
		} else {
			pre[i] = pre[i-1]
		}
		if ch == 'P' {
			p = append(p, i)
		}
	}

	tot := len(p) - 1
	cnt := len(d) - 1

	if tot == 0 {
		fmt.Println("0 0")
		return
	}

	if tot == 1 {
		x := p[1]
		leftStars := pre[x]
		rightStars := cnt - pre[x]

		if leftStars > rightStars {
			fmt.Printf("%d %d\n", leftStars, x-d[1])
		} else if leftStars < rightStars {
			fmt.Printf("%d %d\n", rightStars, d[cnt]-x)
		} else {
			t1 := x - d[1]
			t2 := d[cnt] - x
			if t1 < t2 {
				fmt.Printf("%d %d\n", leftStars, t1)
			} else {
				fmt.Printf("%d %d\n", leftStars, t2)
			}
		}
		return
	}

	empty := func(l, r int) bool {
		if l > r {
			return true
		}
		if l < 1 {
			l = 1
		}
		if r > n {
			r = n
		}
		return pre[l-1] == pre[r]
	}

	check := func(t int) bool {
		f := make([]int, tot+1)
		for i := 1; i <= tot; i++ {
			x := p[i]
			f[i] = 0
			if empty(f[i-1]+1, x-1) {
				if f[i] < x+t {
					f[i] = x + t
				}
			}
			if empty(f[i-1]+1, x-t-1) {
				if f[i] < x {
					f[i] = x
				}
			}
			if i > 1 && empty(f[i-2]+1, x-t-1) {
				if f[i] < p[i-1]+t {
					f[i] = p[i-1] + t
				}
			}
		}
		return f[tot] >= d[cnt]
	}

	l := 0
	r := n
	ans := 0
	for l <= r {
		mid := (l + r) / 2
		if check(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	fmt.Printf("%d %d\n", cnt, ans)
}

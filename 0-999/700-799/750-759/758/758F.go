package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func powLimit(base, exp int, limit int64) int64 {
	res := int64(1)
	b := int64(base)
	for i := 0; i < exp; i++ {
		if res > limit/b {
			return limit + 1
		}
		res *= b
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var l, r int64
	if _, err := fmt.Fscan(reader, &n, &l, &r); err != nil {
		return
	}

	if n == 1 {
		fmt.Println(r - l + 1)
		return
	}
	if n == 2 {
		N := r - l + 1
		fmt.Println(N * (N - 1))
		return
	}
	if n >= 25 {
		fmt.Println(0)
		return
	}

	limit := 1
	for powLimit(limit+1, n-1, r) <= r {
		limit++
	}

	pow := make([]int64, limit+1)
	for i := 1; i <= limit; i++ {
		pow[i] = powLimit(i, n-1, r)
	}

	ans := int64(0)
	for p := 1; p <= limit; p++ {
		pp := pow[p]
		for q := 1; q <= limit; q++ {
			if p == q {
				continue
			}
			if gcd(p, q) != 1 {
				continue
			}
			qq := pow[q]
			tmin := (l + pp - 1) / pp
			if x := (l + qq - 1) / qq; x > tmin {
				tmin = x
			}
			tmax := r / pp
			if x := r / qq; x < tmax {
				tmax = x
			}
			if tmin <= tmax {
				ans += tmax - tmin + 1
			}
		}
	}
	fmt.Println(ans)
}

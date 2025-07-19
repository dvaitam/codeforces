package main

import (
	"fmt"
)

func isPowerOfTwo(x int) bool {
	return x > 0 && x&(x-1) == 0
}

func log2(x int) int {
	k := 0
	for x > 1 {
		x >>= 1
		k++
	}
	return k
}

func cal(x, f int) int {
	if isPowerOfTwo(x) {
		return log2(x) + f
	}
	return log2(x) + 2
}

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	m := n
	ps := make([]int, 0)
	cnt := make([]int, 0)
	for i := 2; i <= m; i++ {
		if m%i == 0 {
			c := 0
			for m%i == 0 {
				m /= i
				c++
			}
			ps = append(ps, i)
			cnt = append(cnt, c)
		}
	}
	if m > 1 {
		ps = append(ps, m)
		cnt = append(cnt, 1)
	}
	res1 := 1
	ma := 1
	f := 0
	for i := 0; i < len(ps); i++ {
		res1 *= ps[i]
		if cnt[i] > ma {
			ma = cnt[i]
		}
		if i+1 < len(cnt) && cnt[i] != cnt[i+1] {
			f = 1
		}
	}
	ops := cal(ma, f)
	fmt.Println(res1, ops)
}

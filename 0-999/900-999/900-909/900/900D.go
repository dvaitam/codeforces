package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

var memo = map[int64]int64{}

func phi(n int64) int64 {
	if val, ok := memo[n]; ok {
		return val
	}
	res := modPow(2, n-1)
	for d := int64(1); d*d <= n; d++ {
		if n%d == 0 {
			if d < n {
				res -= phi(d)
			}
			other := n / d
			if other != d && other < n {
				res -= phi(other)
			}
		}
	}
	res %= mod
	if res < 0 {
		res += mod
	}
	memo[n] = res
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x, y int64
	if _, err := fmt.Fscan(reader, &x, &y); err != nil {
		return
	}
	if y%x != 0 {
		fmt.Println(0)
		return
	}
	m := y / x
	ans := phi(m)
	fmt.Println(ans)
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n   int
	m   int64
	arr []int64
	phi []int64
)

func phiFunc(x int64) int64 {
	res := x
	i := int64(2)
	for i*i <= x {
		if x%i == 0 {
			res = res / i * (i - 1)
			for x%i == 0 {
				x /= i
			}
		}
		i++
	}
	if x > 1 {
		res = res / x * (x - 1)
	}
	return res
}

func buildPhiChain(m int64) []int64 {
	chain := []int64{m}
	for m > 1 {
		m = phiFunc(m)
		chain = append(chain, m)
	}
	if chain[len(chain)-1] != 1 {
		chain = append(chain, 1)
	}
	return chain
}

func powMod(a, e, mod int64) int64 {
	if mod == 1 {
		return 0
	}
	res := int64(1 % mod)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = (res * a) % mod
		}
		a = (a * a) % mod
		e >>= 1
	}
	return res
}

func powOverflow(a, e, limit int64) bool {
	if limit <= 1 {
		return true
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			if a == 0 {
				return false
			}
			if res >= limit || res > (limit-1)/a {
				return true
			}
			res *= a
		}
		e >>= 1
		if e == 0 {
			break
		}
		if a >= limit || a > (limit-1)/a {
			return true
		}
		a *= a
	}
	return res >= limit
}

func calc(l, r, idx int) (int64, bool) {
	mod := phi[idx]
	if mod == 1 {
		return 0, true
	}
	if l == r {
		val := arr[l] % mod
		big := arr[l] >= mod
		return val, big
	}
	valNext, bigNext := calc(l+1, r, idx+1)
	exp := valNext
	if bigNext {
		exp += phi[idx+1]
	}
	val := powMod(arr[l], exp, mod)
	big := powOverflow(arr[l], exp, mod)
	return val, big
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &m)
	arr = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	phi = buildPhiChain(m)
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		res, _ := calc(l, r, 0)
		fmt.Fprintln(writer, res%int64(m))
	}
}

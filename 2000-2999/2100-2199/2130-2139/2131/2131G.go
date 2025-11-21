package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	mod = 1_000_000_007
	inf = int64(1e18)
)

var pow2 [65]int64

func init() {
	pow2[0] = 1
	for i := 1; i < len(pow2); i++ {
		pow2[i] = pow2[i-1] << 1
		if pow2[i] > inf {
			pow2[i] = inf
		}
	}
}

func lenV(m int64) int64 {
	if m <= 0 {
		return 0
	}
	if m >= 60 {
		return inf
	}
	return pow2[m] - 1
}

func prodV(m, L int64) int64 {
	if L <= 0 || m <= 0 {
		return 1
	}
	origM := m
	if m > 60 {
		m = 60
	}
	lenPrev := lenV(m - 1)
	if L <= lenPrev {
		return prodV(origM-1, L)
	}
	leftProd := prodV(origM-1, lenPrev)
	res := leftProd * (origM % mod) % mod
	if L == lenPrev+1 {
		return res
	}
	rightLen := L - lenPrev - 1
	return res * prodV(origM-1, rightLen) % mod
}

func prodT(x, L int64) int64 {
	if L == 0 {
		return 1
	}
	res := x % mod
	if L == 1 {
		return res
	}
	return res * prodV(x-1, L-1) % mod
}

func lenT(x int64) int64 {
	if x <= 1 {
		return 1
	}
	m := x - 1
	if m >= 60 {
		return inf
	}
	return pow2[m]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
		}
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

		ans := int64(1)
		remaining := k
		for _, v := range vals {
			if remaining == 0 {
				break
			}
			length := lenT(v)
			take := remaining
			if length < take {
				take = length
			}
			ans = ans * prodT(v, take) % mod
			remaining -= take
		}
		fmt.Fprintln(out, ans%mod)
	}
}

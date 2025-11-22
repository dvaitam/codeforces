package main

import (
	"bufio"
	"fmt"
	"os"
)

func requiredDecrements(a, b []int64, pre []int64, bsum, amax int64, n int64, T int64) int64 {
	if T == 0 {
		var need int64
		for _, v := range a {
			if v > 0 {
				need += v
			}
		}
		return need
	}
	q := T / n
	r := int(T % n)

	limitQ := amax/bsum + 2
	var base int64
	if q >= limitQ {
		base = amax + bsum // definitely enough to avoid deficits overflow
	} else {
		base = q * bsum
	}

	var need int64
	for i := 0; i < len(a); i++ {
		partial := pre[i+r] - pre[i]
		s := base + partial
		if s < a[i] {
			need += a[i] - s
			if need < 0 { // overflow guard (unlikely)
				return 1 << 62
			}
		}
	}
	return need
}

func solveCase(n int, k int64, a, b []int64) int64 {
	var bsum, amax int64
	for i := 0; i < n; i++ {
		bsum += b[i]
		if a[i] > amax {
			amax = a[i]
		}
	}

	dup := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		dup[i+1] = dup[i] + b[i%n]
	}

	need0 := requiredDecrements(a, b, dup, bsum, amax, int64(n), 0)
	if need0 <= k {
		return 0
	}

	high := int64(1)
	for requiredDecrements(a, b, dup, bsum, amax, int64(n), high) > k {
		high <<= 1
	}
	low := high >> 1

	for high-low > 1 {
		mid := (low + high) / 2
		if requiredDecrements(a, b, dup, bsum, amax, int64(n), mid) <= k {
			high = mid
		} else {
			low = mid
		}
	}
	return high
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		ans := solveCase(n, k, a, b)
		fmt.Fprintln(out, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 200000
const inf = int(1e9)

func powInt(a int64, b int) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res *= a
		}
		a *= a
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	spf := make([]int, maxA+1)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	min1 := make([]int, maxA+1)
	min2 := make([]int, maxA+1)
	cnt := make([]int, maxA+1)
	for i := 0; i <= maxA; i++ {
		min1[i] = inf
		min2[i] = inf
	}

	for _, x0 := range a {
		x := x0
		for x > 1 {
			p := spf[x]
			e := 0
			for x%p == 0 {
				x /= p
				e++
			}
			cnt[p]++
			if e < min1[p] {
				min2[p] = min1[p]
				min1[p] = e
			} else if e < min2[p] {
				min2[p] = e
			}
		}
	}

	var ans int64 = 1
	for p := 2; p <= maxA; p++ {
		if cnt[p] == 0 {
			continue
		}
		notDiv := n - cnt[p]
		var exp int
		if notDiv >= 2 {
			continue
		} else if notDiv == 1 {
			exp = min1[p]
		} else {
			exp = min2[p]
		}
		if exp == inf {
			continue
		}
		ans *= powInt(int64(p), exp)
	}

	fmt.Fprintln(out, ans)
}

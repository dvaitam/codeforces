package main

import (
	"bufio"
	"fmt"
	"os"
)

// minSum returns the minimal possible sum of the sequence when a_n is fixed.
// It stops early if the sum exceeds limit to avoid overflow.
func minSum(n, an, limit int64) int64 {
	sum := an
	cur := an
	if sum > limit {
		return limit + 1
	}
	for i := int64(1); i < n; i++ {
		if cur > 1 {
			cur = (cur + 1) / 2
		}
		if sum > limit-cur {
			return limit + 1
		}
		sum += cur
		if cur == 1 {
			rest := n - 1 - i
			if rest > 0 {
				if sum > limit-rest {
					return limit + 1
				}
				sum += rest
			}
			break
		}
	}
	return sum
}

// solve returns the maximal possible value of a_n for given n and s.
func solve(n, s int64) int64 {
	l, r := int64(1), s
	ans := int64(1)
	for l <= r {
		mid := (l + r) / 2
		if minSum(n, mid, s) <= s {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, s int64
		fmt.Fscan(in, &n, &s)
		fmt.Fprintln(out, solve(n, s))
	}
}

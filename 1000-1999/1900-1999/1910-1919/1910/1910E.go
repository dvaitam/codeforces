package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		M := make([]int64, n)
		m := make([]int64, n)
		for i := 0; i < n; i++ {
			if a[i] > b[i] {
				M[i] = a[i]
				m[i] = b[i]
			} else {
				M[i] = b[i]
				m[i] = a[i]
			}
		}

		left := make([]int64, n)
		cur := int64(0)
		best := int64(0)
		for i := 0; i < n; i++ {
			cur += M[i]
			if cur < 0 {
				cur = 0
			}
			if cur > best {
				best = cur
			}
			left[i] = best
		}
		right := make([]int64, n)
		cur = 0
		best = 0
		for i := n - 1; i >= 0; i-- {
			cur += M[i]
			if cur < 0 {
				cur = 0
			}
			if cur > best {
				best = cur
			}
			right[i] = best
		}
		bestTwo := left[n-1]
		for i := 0; i < n-1; i++ {
			val := left[i] + right[i+1]
			if val > bestTwo {
				bestTwo = val
			}
		}

		prefM := make([]int64, n+1)
		prefm := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefM[i+1] = prefM[i] + M[i]
			prefm[i+1] = prefm[i] + m[i]
		}
		minPrefM := make([]int64, n+1)
		mn := int64(0)
		for i := 1; i <= n; i++ {
			if prefM[i] < mn {
				mn = prefM[i]
			}
			minPrefM[i] = mn
		}
		minComb := make([]int64, n+1)
		val := int64(1 << 60)
		for t1 := 1; t1 <= n; t1++ {
			c := prefm[t1-1] + minPrefM[t1-1]
			if c < val {
				val = c
			}
			minComb[t1] = val
		}
		maxPrefFrom := make([]int64, n+1)
		mx := int64(-1 << 60)
		for i := n; i >= 0; i-- {
			if prefM[i] > mx {
				mx = prefM[i]
			}
			maxPrefFrom[i] = mx
		}
		bestOverlap := int64(0)
		for t1 := 1; t1 <= n; t1++ {
			cand := prefm[t1] + maxPrefFrom[t1] - minComb[t1]
			if cand > bestOverlap {
				bestOverlap = cand
			}
		}
		ans := bestTwo
		if bestOverlap > ans {
			ans = bestOverlap
		}
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}

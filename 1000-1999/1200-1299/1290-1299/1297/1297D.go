package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int64
	idx int
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
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		p := make([]pair, n)
		for i := 0; i < n; i++ {
			p[i] = pair{arr[i], i}
		}
		sort.Slice(p, func(i, j int) bool { return p[i].val < p[j].val })
		sorted := make([]int64, n)
		for i := 0; i < n; i++ {
			sorted[i] = p[i].val
		}
		baseline := int64(0)
		for i := 0; i < n; i++ {
			v := sorted[i] + int64(n-1-i)
			if v > baseline {
				baseline = v
			}
		}
		lo := baseline
		hi := baseline + k
		for lo < hi {
			mid := (lo + hi) / 2
			var cap int64
			for i := 0; i < n; i++ {
				allowed := mid - int64(n-1-i)
				if allowed > sorted[i] {
					cap += allowed - sorted[i]
					if cap >= k {
						break
					}
				}
			}
			if cap >= k {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		M := lo
		res := make([]int64, n)
		remaining := k
		for i := n - 1; i >= 0; i-- {
			allowed := M - int64(n-1-i)
			if i != n-1 {
				next := res[i+1] - 1
				if allowed > next {
					allowed = next
				}
			}
			if allowed < sorted[i] {
				allowed = sorted[i]
			}
			inc := allowed - sorted[i]
			if inc > remaining {
				inc = remaining
			}
			res[i] = sorted[i] + inc
			remaining -= inc
		}
		if remaining > 0 {
			for i := n - 1; i >= 0 && remaining > 0; i-- {
				allowed := M - int64(n-1-i)
				if i != n-1 {
					next := res[i+1] - 1
					if allowed > next {
						allowed = next
					}
				}
				if allowed <= res[i] {
					continue
				}
				inc := allowed - res[i]
				if inc > remaining {
					inc = remaining
				}
				res[i] += inc
				remaining -= inc
			}
		}
		bonus := make([]int64, n)
		for i := 0; i < n; i++ {
			bonus[p[i].idx] = res[i] - p[i].val
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, bonus[i])
		}
		fmt.Fprintln(out)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1e9 + 7
const N int64 = 1e18

type Interval struct {
	L, R int64
	g    int64
	sum  int64
}

func powLimit(base, exp, limit int64) int64 {
	res := int64(1)
	for ; exp > 0; exp-- {
		if base != 0 && res > limit/base {
			return limit + 1
		}
		res *= base
	}
	return res
}

func generateIntervals() []Interval {
	var intervals []Interval
	for k := int64(2); ; k++ {
		L := int64(1) << k
		if L > N {
			break
		}
		var R int64
		if k+1 >= 63 {
			R = N
		} else {
			R = (int64(1) << (k + 1)) - 1
			if R > N {
				R = N
			}
		}
		for z := int64(1); ; z++ {
			pz := powLimit(k, z, N)
			if pz > R {
				break
			}
			start := L
			if pz > start {
				start = pz
			}
			pn := powLimit(k, z+1, N)
			end := R
			if pn-1 < end {
				end = pn - 1
			}
			if start <= end {
				intervals = append(intervals, Interval{L: start, R: end, g: z})
			}
			if pn > R {
				break
			}
		}
	}
	var prefix int64
	for i := range intervals {
		length := intervals[i].R - intervals[i].L + 1
		prefix = (prefix + (intervals[i].g%mod)*(length%mod)) % mod
		intervals[i].sum = prefix
	}
	return intervals
}

func prefixSum(intervals []Interval, x int64) int64 {
	if x < 4 {
		return 0
	}
	idx := sort.Search(len(intervals), func(i int) bool { return intervals[i].R >= x })
	if idx == len(intervals) {
		return intervals[len(intervals)-1].sum
	}
	var res int64
	if idx > 0 {
		res = intervals[idx-1].sum
	}
	if x >= intervals[idx].L {
		length := x - intervals[idx].L + 1
		res = (res + (intervals[idx].g%mod)*(length%mod)) % mod
	}
	return res
}

func main() {
	intervals := generateIntervals()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		ans := prefixSum(intervals, r) - prefixSum(intervals, l-1)
		ans %= mod
		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(out, ans)
	}
}

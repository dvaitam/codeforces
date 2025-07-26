package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for problemB.txt from contest 1801.
// We choose exactly one store in each department so that both friends get at
// least one gift while minimizing the difference between the most expensive
// gifts they receive. Fixing a department for friend1 with price a_i = m1,
// every department with a_j > m1 must go to friend2 which forces m2 to be at
// least the maximum b_j among them. If there are none, friend2 must get a gift
// from some other department so m2 is at least the minimum b_j excluding this
// department. Among all prices not from the same department and not smaller than
// this threshold we pick the one closest to m1. The same procedure is repeated
// swapping the roles of a and b. Sorting allows an O(n log n) solution.

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(n int, a, b []int) int {
	type pair struct{ val, oth, idx int }

	compute := func(val, oth []int) int {
		pairs := make([]pair, n)
		for i := 0; i < n; i++ {
			pairs[i] = pair{val[i], oth[i], i}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
		vs := make([]int, n)
		os := make([]int, n)
		idx := make([]int, n)
		for i, p := range pairs {
			vs[i] = p.val
			os[i] = p.oth
			idx[i] = p.idx
		}
		sufMax := make([]int, n)
		maxVal := 0
		for i := n - 1; i >= 0; {
			j := i
			for j >= 0 && vs[j] == vs[i] {
				sufMax[j] = maxVal
				j--
			}
			g := 0
			for k := j + 1; k <= i; k++ {
				if os[k] > g {
					g = os[k]
				}
			}
			if g > maxVal {
				maxVal = g
			}
			i = j
		}
		preMin := make([]int, n)
		sufMin := make([]int, n)
		for i := 0; i < n; i++ {
			if i == 0 || os[i] < preMin[i-1] {
				preMin[i] = os[i]
			} else {
				preMin[i] = preMin[i-1]
			}
		}
		for i := n - 1; i >= 0; i-- {
			if i == n-1 || os[i] < sufMin[i+1] {
				sufMin[i] = os[i]
			} else {
				sufMin[i] = sufMin[i+1]
			}
		}
		minExcl := make([]int, n)
		const inf = int(2e9 + 5)
		for i := 0; i < n; i++ {
			left, right := inf, inf
			if i > 0 {
				left = preMin[i-1]
			}
			if i+1 < n {
				right = sufMin[i+1]
			}
			if left < right {
				minExcl[i] = left
			} else {
				minExcl[i] = right
			}
		}
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{oth[i], val[i], i}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		vals := make([]int, n)
		for i, p := range arr {
			vals[i] = p.val
		}
		best := inf
		for pos := 0; pos < n; pos++ {
			id := idx[pos]
			v := vs[pos]
			thr := sufMax[pos]
			if thr == 0 {
				thr = minExcl[pos]
			}
			if thr >= inf {
				continue
			}
			lo := v
			if thr > lo {
				lo = thr
			}
			p := sort.Search(n, func(i int) bool { return vals[i] >= lo })
			for _, t := range []int{p - 1, p, p + 1} {
				if t >= 0 && t < n {
					cand := arr[t]
					if cand.idx == id {
						continue
					}
					if cand.val >= thr {
						diff := abs(v - cand.val)
						if diff < best {
							best = diff
						}
					}
				}
			}
		}
		return best
	}

	ans := compute(a, b)
	t := compute(b, a)
	if t < ans {
		ans = t
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}
		fmt.Fprintln(out, solveCase(n, a, b))
	}
}

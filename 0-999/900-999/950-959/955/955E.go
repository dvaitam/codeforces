package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct{ l, r int }

func canTrap(a []int, t int) bool {
	n := len(a)
	if t >= n { // need at least one icicle to the right
		return false
	}
	left := make([]interval, 0)
	right := make([]interval, 0)
	for i := 1; i <= n; i++ {
		need := t - (a[i-1] - 1)
		if need < 0 {
			continue
		}
		l := i - need
		r := i + need
		if l < 1 {
			l = 1
		}
		if r > n {
			r = n
		}
		if i <= t {
			left = append(left, interval{l, r})
		} else {
			right = append(right, interval{l, r})
		}
	}
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	sort.Slice(left, func(i, j int) bool { return left[i].l < left[j].l })
	sort.Slice(right, func(i, j int) bool { return right[i].l < right[j].l })
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i].r < right[j].l {
			i++
		} else if right[j].r < left[i].l {
			j++
		} else {
			return true
		}
	}
	return false
}

func minTime(a []int) int {
	n := len(a)
	maxA := 0
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	lo, hi := 1, n+maxA
	ans := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if canTrap(a, mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	fmt.Println(minTime(a))
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solve(n, x, y, a))
	}
}

func solve(n int, x, y int64, a []int64) int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		v := a[i-1] - int64(i-1)
		if v > pref[i-1] {
			pref[i] = v
		} else {
			pref[i] = pref[i-1]
		}
	}
	if x < pref[1] {
		return -1
	}
	r := x
	var games int64
	for r < y {
		k := sort.Search(len(pref), func(i int) bool { return pref[i] > r }) - 1
		if r+int64(k) >= y {
			games += y - r
			break
		}
		delta := int64(2*k - n)
		if delta <= 0 {
			return -1
		}
		if k == n {
			games += y - r
			break
		}
		nextR := pref[k+1]
		finish := y - int64(k)
		target := nextR
		if finish < target {
			target = finish
		}
		cycles := (target - r + delta - 1) / delta
		if cycles <= 0 {
			cycles = 1
		}
		r += cycles * delta
		games += cycles * int64(n)
	}
	return games
}

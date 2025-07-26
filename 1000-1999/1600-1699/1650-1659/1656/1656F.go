package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func cost(a []int64, pref []int64, t int64) int64 {
	n := len(a)
	aMin := a[0]
	aMax := a[n-1]
	idx := sort.Search(n, func(i int) bool { return a[i] >= -t })
	sumLess := pref[idx]
	cntLess := int64(idx)
	sumGreater := pref[n] - pref[idx]
	cntGreater := int64(n - idx)
	if idx > 0 {
		sumLess -= aMin
		cntLess--
	} else {
		sumGreater -= aMin
		cntGreater--
	}
	if n-1 < idx {
		sumLess -= aMax
		cntLess--
	} else {
		sumGreater -= aMax
		cntGreater--
	}
	res := aMin*aMax + t*(aMin+aMax)
	res += aMin*sumGreater + t*sumGreater + cntGreater*t*aMin
	res += aMax*sumLess + t*sumLess + cntLess*t*aMax
	return res
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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + a[i]
		}
		sumAll := pref[n]
		slopePlus := sumAll + int64(n-2)*a[0]
		slopeMinus := sumAll + int64(n-2)*a[n-1]
		if slopePlus > 0 || slopeMinus < 0 {
			fmt.Fprintln(out, "INF")
			continue
		}
		candidates := make(map[int64]struct{})
		for i := 0; i < n; i++ {
			candidates[-a[i]] = struct{}{}
		}
		for i := 0; i < n-1; i++ {
			candidates[-a[n-1]-a[i]] = struct{}{}
		}
		for i := 1; i < n; i++ {
			candidates[-a[0]-a[i]] = struct{}{}
		}
		best := int64(-1 << 63)
		for t := range candidates {
			v := cost(a, pref, t)
			if v > best {
				best = v
			}
		}
		fmt.Fprintln(out, best)
	}
}

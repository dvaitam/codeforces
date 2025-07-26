package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func divisors(n int64) []int64 {
	res := []int64{}
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			res = append(res, i)
			if i*i != n {
				res = append(res, n/i)
			}
		}
	}
	return res
}

func valid(a, w, l int64) bool {
	for mask := 0; mask < 16; mask++ {
		tl := int64((mask >> 0) & 1)
		tr := int64((mask >> 1) & 1)
		br := int64((mask >> 2) & 1)
		bl := int64((mask >> 3) & 1)
		top := w - tl - tr
		bottom := w - bl - br
		left := l - (1 - tl) - (1 - bl)
		right := l - (1 - tr) - (1 - br)
		if top%a == 0 && bottom%a == 0 && left%a == 0 && right%a == 0 {
			return true
		}
	}
	return false
}

func solve(w, l int64) []int64 {
	candidate := map[int64]struct{}{}
	for _, v := range []int64{w, w - 1, w - 2, l, l - 1, l - 2} {
		if v <= 0 {
			continue
		}
		for _, d := range divisors(v) {
			candidate[d] = struct{}{}
		}
	}
	list := make([]int64, 0, len(candidate))
	for v := range candidate {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	ans := make([]int64, 0, len(list))
	for _, a := range list {
		if valid(a, w, l) {
			ans = append(ans, a)
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var w, l int64
		fmt.Fscan(in, &w, &l)
		res := solve(w, l)
		fmt.Print(len(res))
		for _, v := range res {
			fmt.Print(" ", v)
		}
		fmt.Println()
	}
}

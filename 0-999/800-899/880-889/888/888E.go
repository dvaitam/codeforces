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

	var n int
	var m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		arr[i] %= m
	}

	mid := n / 2
	left := arr[:mid]
	right := arr[mid:]

	// compute subset sums mod m for a slice
	compute := func(nums []int64) []int64 {
		res := []int64{0}
		for _, v := range nums {
			size := len(res)
			for i := 0; i < size; i++ {
				res = append(res, (res[i]+v)%m)
			}
		}
		return res
	}

	sums1 := compute(left)
	sums2 := compute(right)
	sort.Slice(sums2, func(i, j int) bool { return sums2[i] < sums2[j] })

	var ans int64
	for _, x := range sums1 {
		// best y <= m-1-x
		target := m - 1 - x
		idx := sort.Search(len(sums2), func(i int) bool { return sums2[i] > target })
		if idx > 0 {
			if x+sums2[idx-1] > ans {
				ans = x + sums2[idx-1]
			}
		}
		cand := (x + sums2[len(sums2)-1]) % m
		if cand > ans {
			ans = cand
		}
	}

	fmt.Fprintln(out, ans)
}

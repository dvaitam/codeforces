package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func upperBound(a []int64, target int64) int {
	l, r := 0, len(a)
	for l < r {
		mid := (l + r) >> 1
		if a[mid] <= target {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		best := 1
		for l := 0; l < n; l++ {
			if l == n-1 {
				if best < 1 {
					best = 1
				}
				continue
			}
			s := arr[l] + arr[l+1] - 1
			r := upperBound(arr, s) - 1
			if r < l+1 {
				r = l + 1
			}
			length := r - l + 1
			if length > best {
				best = length
			}
		}
		fmt.Fprintln(out, n-best)
	}
}

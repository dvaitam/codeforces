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
		fmt.Fscan(in, &n)
		cnt := make([]int, n)
		for i := 2; i <= n; i++ {
			var p int
			fmt.Fscan(in, &p)
			cnt[p-1]++
		}
		var arr []int
		for _, c := range cnt {
			if c > 0 {
				arr = append(arr, c)
			}
		}
		arr = append(arr, 1)
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		m := len(arr)
		for i := 0; i < m; i++ {
			arr[i] -= m - i
			if arr[i] < 0 {
				arr[i] = 0
			}
		}
		var b []int
		for _, v := range arr {
			if v > 0 {
				b = append(b, v)
			}
		}
		if len(b) == 0 {
			fmt.Fprintln(out, m)
			continue
		}
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		lo, hi := 0, n+5
		for lo < hi {
			mid := (lo + hi) / 2
			sum := 0
			for _, v := range b {
				if v > mid {
					sum += v - mid
				}
			}
			if sum <= mid {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(out, m+lo)
	}
}

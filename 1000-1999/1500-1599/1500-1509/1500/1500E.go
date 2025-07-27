package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// calcUnsuitable computes the number of unsuitable integers for a sorted array.
func calcUnsuitable(arr []int64) int64 {
	n := len(arr)
	if n <= 1 {
		return 0
	}
	// prefix sums
	pre := make([]int64, n+1)
	for i, v := range arr {
		pre[i+1] = pre[i] + v
	}
	total := pre[n]
	base := total - 2*arr[0]
	if n == 2 {
		if base < 0 {
			return 0
		}
		return base
	}
	gaps := int64(0)
	for k := 1; k <= n-2; k++ {
		gap := pre[k+1] + pre[n-k] - total
		if gap > 0 {
			gaps += gap
		}
	}
	res := base - gaps
	if res < 0 {
		return 0
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	fmt.Fprintln(out, calcUnsuitable(arr))
	for ; q > 0; q-- {
		var t int
		var a int64
		fmt.Fscan(in, &t, &a)
		if t == 1 { // add
			arr = append(arr, a)
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		} else { // remove
			// find index
			idx := sort.Search(len(arr), func(i int) bool { return arr[i] >= a })
			if idx < len(arr) && arr[idx] == a {
				arr = append(arr[:idx], arr[idx+1:]...)
			}
		}
		fmt.Fprintln(out, calcUnsuitable(arr))
	}
}

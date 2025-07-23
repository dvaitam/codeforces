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
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	l, r := 0, n-1
	for l < r && k > 0 {
		if l+1 <= n-1-r {
			diff := arr[l+1] - arr[l]
			cost := diff * int64(l+1)
			if cost <= k {
				k -= cost
				l++
			} else {
				arr[l] += k / int64(l+1)
				k = 0
			}
		} else {
			diff := arr[r] - arr[r-1]
			cost := diff * int64(n-r)
			if cost <= k {
				k -= cost
				r--
			} else {
				arr[r] -= k / int64(n-r)
				k = 0
			}
		}
	}

	if arr[r] < arr[l] {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, arr[r]-arr[l])
	}
}

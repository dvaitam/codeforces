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
		var n, k, x int
		fmt.Fscan(in, &n, &k, &x)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + int64(arr[i])
		}
		ans := int64(-1 << 60)
		if k > n {
			k = n
		}
		for j := 0; j <= k; j++ {
			idx := j + x
			if idx > n {
				idx = n
			}
			val := prefix[n] - 2*prefix[idx] + prefix[j]
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}

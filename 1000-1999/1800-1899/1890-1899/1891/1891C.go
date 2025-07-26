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
		arr := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			sum += arr[i]
		}
		limit := sum / 2
		if limit == 0 {
			fmt.Fprintln(out, sum)
			continue
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		var pref int64
		var r int64
		for _, v := range arr {
			pref += v
			r++
			if pref >= limit {
				break
			}
		}
		ans := sum - limit + r
		fmt.Fprintln(out, ans)
	}
}

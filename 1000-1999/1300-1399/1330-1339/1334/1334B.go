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
		var x int64
		fmt.Fscan(in, &n, &x)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		var sum int64
		ans := 0
		for i := 0; i < n; i++ {
			sum += arr[i]
			if sum >= x*int64(i+1) {
				ans = i + 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}

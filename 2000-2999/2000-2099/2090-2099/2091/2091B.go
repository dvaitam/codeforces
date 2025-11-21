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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		var ans int64
		count := int64(0)
		for _, v := range a {
			count++
			if v*count >= x {
				ans++
				count = 0
			}
		}
		fmt.Fprintln(out, ans)
	}
}

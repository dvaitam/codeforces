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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Ints(a)
		maxLen := 1
		cur := 1
		for i := 1; i < n; i++ {
			if a[i]-a[i-1] <= k {
				cur++
			} else {
				if cur > maxLen {
					maxLen = cur
				}
				cur = 1
			}
		}
		if cur > maxLen {
			maxLen = cur
		}
		fmt.Fprintln(out, n-maxLen)
	}
}

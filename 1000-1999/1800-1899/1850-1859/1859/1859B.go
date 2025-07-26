package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt from contest 1859.
// For each array we only benefit from removing its smallest value,
// raising its minimum to the second smallest. All removed values can
// be moved to a single array whose minimum then equals the global
// smallest element. Thus the maximum achievable beauty is:
//
//	minFirst + sum(secondMin) - min(secondMin).
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		minFirst := int64(1 << 62)
		minSecond := int64(1 << 62)
		var sumSecond int64
		for i := 0; i < n; i++ {
			var m int
			fmt.Fscan(in, &m)
			first, second := int64(1<<62), int64(1<<62)
			for j := 0; j < m; j++ {
				var x int64
				fmt.Fscan(in, &x)
				if x < first {
					second = first
					first = x
				} else if x < second {
					second = x
				}
			}
			if first < minFirst {
				minFirst = first
			}
			if second < minSecond {
				minSecond = second
			}
			sumSecond += second
		}
		ans := minFirst + sumSecond - minSecond
		fmt.Fprintln(out, ans)
	}
}

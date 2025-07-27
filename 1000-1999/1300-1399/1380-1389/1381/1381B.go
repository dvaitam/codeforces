package main

import (
	"bufio"
	"fmt"
	"os"
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
		size := 2 * n
		p := make([]int, size)
		for i := 0; i < size; i++ {
			fmt.Fscan(in, &p[i])
		}
		segments := make([]int, 0)
		curMax := p[0]
		segLen := 1
		for i := 1; i < size; i++ {
			if p[i] > curMax {
				segments = append(segments, segLen)
				segLen = 1
				curMax = p[i]
			} else {
				segLen++
			}
		}
		segments = append(segments, segLen)

		dp := make([]bool, n+1)
		dp[0] = true
		for _, l := range segments {
			for i := n; i >= l; i-- {
				if dp[i-l] {
					dp[i] = true
				}
			}
		}
		if dp[n] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		nextBad := make([]int, n+2)
		for i := range nextBad {
			nextBad[i] = n + 1
		}
		for i := 0; i < m; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x > y {
				x, y = y, x
			}
			if y < nextBad[x] {
				nextBad[x] = y
			}
		}
		minVal := n + 1
		var ans int64
		for l := n; l >= 1; l-- {
			if nextBad[l] < minVal {
				minVal = nextBad[l]
			}
			limit := minVal - 1
			if limit > n {
				limit = n
			}
			ans += int64(limit - l + 1)
		}
		fmt.Fprintln(writer, ans)
	}
}

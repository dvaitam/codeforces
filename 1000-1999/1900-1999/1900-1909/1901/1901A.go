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
		var n, x int
		fmt.Fscan(reader, &n, &x)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		maxGap := a[0]
		for i := 1; i < n; i++ {
			gap := a[i] - a[i-1]
			if gap > maxGap {
				maxGap = gap
			}
		}
		lastGap := 2 * (x - a[n-1])
		if lastGap > maxGap {
			maxGap = lastGap
		}
		fmt.Fprintln(writer, maxGap)
	}
}

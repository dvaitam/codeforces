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

	var q int
	fmt.Fscan(reader, &q)
	for q > 0 {
		q--
		var n, t int
		fmt.Fscan(reader, &n, &t)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		bestIdx := -1
		bestVal := -1
		for i := 0; i < n; i++ {
			if a[i]+i <= t && b[i] > bestVal {
				bestVal = b[i]
				bestIdx = i + 1
			}
		}
		fmt.Fprintln(writer, bestIdx)
	}
}

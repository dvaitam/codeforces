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
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// Check if already not sorted
		notSorted := false
		for i := 0; i < n-1; i++ {
			if a[i] > a[i+1] {
				notSorted = true
				break
			}
		}
		if notSorted {
			fmt.Fprintln(writer, 0)
			continue
		}
		// Array is non-decreasing
		minDelta := int(1<<31 - 1)
		for i := 0; i < n-1; i++ {
			delta := a[i+1] - a[i]
			if delta < minDelta {
				minDelta = delta
			}
		}
		operations := minDelta/2 + 1
		fmt.Fprintln(writer, operations)
	}
}

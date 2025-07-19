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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// find position of maximum element (last occurrence if tie)
		pos := 1
		maxVal := a[0]
		for i := 1; i < n; i++ {
			if a[i] > maxVal || (a[i] == maxVal && i+1 > pos) {
				maxVal = a[i]
				pos = i + 1
			}
		}
		fmt.Fprintln(writer, pos)
	}
}

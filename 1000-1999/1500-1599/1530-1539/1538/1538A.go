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
		minIdx, maxIdx := 0, 0
		for i := 1; i < n; i++ {
			if a[i] < a[minIdx] {
				minIdx = i
			}
			if a[i] > a[maxIdx] {
				maxIdx = i
			}
		}
		// convert to 1-based positions
		minPos := minIdx + 1
		maxPos := maxIdx + 1
		if minPos > maxPos {
			minPos, maxPos = maxPos, minPos
		}
		fromLeft := maxPos
		fromRight := n - minPos + 1
		mixed1 := minPos + (n - maxPos + 1)
		mixed2 := maxPos + (n - minPos + 1)
		ans := fromLeft
		if fromRight < ans {
			ans = fromRight
		}
		if mixed1 < ans {
			ans = mixed1
		}
		if mixed2 < ans {
			ans = mixed2
		}
		fmt.Fprintln(writer, ans)
	}
}

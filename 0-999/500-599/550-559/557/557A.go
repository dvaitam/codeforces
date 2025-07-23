package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var min1, max1, min2, max2, min3, max3 int
	if _, err := fmt.Fscan(reader, &min1, &max1); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &min2, &max2); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &min3, &max3); err != nil {
		return
	}

	x1 := min(max1, n-min2-min3)
	if x1 < min1 {
		x1 = min1
	}
	remaining := n - x1
	x2 := min(max2, remaining-min3)
	if x2 < min2 {
		x2 = min2
	}
	x3 := n - x1 - x2
	fmt.Fprintln(writer, x1, x2, x3)
}

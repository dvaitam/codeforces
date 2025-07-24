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

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	if n <= m {
		fmt.Fprintln(writer, n)
		return
	}

	if n <= m+1 {
		fmt.Fprintln(writer, n)
		return
	}

	remaining := n - (m + 1)
	left, right := int64(0), int64(2_000_000_000)
	for left < right {
		mid := (left + right) / 2
		if mid*(mid+3) >= 2*remaining {
			right = mid
		} else {
			left = mid + 1
		}
	}
	fmt.Fprintln(writer, m+1+left)
}

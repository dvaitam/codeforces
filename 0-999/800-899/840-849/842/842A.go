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

	var l, r, x, y, k int64
	if _, err := fmt.Fscan(reader, &l, &r, &x, &y, &k); err != nil {
		return
	}
	low := (l + k - 1) / k
	high := r / k
	if low <= high && maxInt64(low, x) <= minInt64(high, y) {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

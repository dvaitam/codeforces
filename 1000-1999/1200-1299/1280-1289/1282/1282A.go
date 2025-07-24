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
		var a, b, c, r int64
		fmt.Fscan(reader, &a, &b, &c, &r)
		if a > b {
			a, b = b, a
		}
		left := c - r
		right := c + r
		start := max(a, left)
		end := min(b, right)
		inter := int64(0)
		if end > start {
			inter = end - start
		}
		fmt.Fprintln(writer, b-a-inter)
	}
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

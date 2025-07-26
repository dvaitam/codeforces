package main

import (
	"bufio"
	"fmt"
	"os"
)

func minLines(n int64) int64 {
	l, r := int64(0), int64(1)
	for r*(r-1) < n {
		r <<= 1
	}
	for l < r {
		m := (l + r) / 2
		if m*(m-1) >= n {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, minLines(n))
	}
}

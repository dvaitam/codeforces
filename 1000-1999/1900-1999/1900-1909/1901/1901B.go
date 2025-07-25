package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
	}
	var segments int64
	var prev int64
	for i := 0; i < n; i++ {
		if c[i] > prev {
			segments += c[i] - prev
		}
		prev = c[i]
	}
	fmt.Fprintln(writer, segments-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for i := 0; i < t; i++ {
		solve(reader, writer)
	}
}

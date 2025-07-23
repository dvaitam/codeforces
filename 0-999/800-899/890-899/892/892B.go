package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	l := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i])
	}
	left := n + 1
	alive := 0
	for i := n; i >= 1; i-- {
		if i < left {
			alive++
		}
		if t := i - l[i-1]; t < left {
			left = t
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, alive)
	writer.Flush()
}

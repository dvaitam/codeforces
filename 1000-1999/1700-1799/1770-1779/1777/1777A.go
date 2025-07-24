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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		runs := 1
		prev := arr[0] % 2
		for i := 1; i < n; i++ {
			cur := arr[i] % 2
			if cur != prev {
				runs++
				prev = cur
			}
		}
		ops := n - runs
		fmt.Fprintln(writer, ops)
	}
}

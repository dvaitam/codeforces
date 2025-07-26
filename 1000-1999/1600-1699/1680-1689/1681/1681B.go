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
		var m int
		fmt.Fscan(reader, &m)
		shift := 0
		for i := 0; i < m; i++ {
			var b int
			fmt.Fscan(reader, &b)
			shift = (shift + b) % n
		}
		fmt.Fprintln(writer, a[shift])
	}
}

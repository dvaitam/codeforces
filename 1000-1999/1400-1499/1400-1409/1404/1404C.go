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

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		l := x + 1
		r := n - y
		rem := 0
		for i := l; i <= r; i++ {
			if a[i-1] <= i && rem >= i-a[i-1] {
				rem++
			}
		}
		fmt.Fprintln(writer, rem)
	}
}

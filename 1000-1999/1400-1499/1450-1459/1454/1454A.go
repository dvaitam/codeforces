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
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		if n == 2 {
			fmt.Fprintln(writer, "2 1")
			continue
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = n - i
		}
		mid := n / 2
		if mid < n-1 {
			a[mid], a[mid+1] = a[mid+1], a[mid]
		}
		for i, v := range a {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
}

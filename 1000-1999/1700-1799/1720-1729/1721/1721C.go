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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		mn := make([]int, n)
		j := 0
		for i := 0; i < n; i++ {
			for j < n && b[j] < a[i] {
				j++
			}
			mn[i] = b[j] - a[i]
		}
		mx := make([]int, n)
		pos := n - 1
		for i := n - 1; i >= 0; i-- {
			mx[i] = b[pos] - a[i]
			if i > 0 && b[i-1] < a[i] {
				pos = i - 1
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, mn[i])
		}
		writer.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, mx[i])
		}
		writer.WriteByte('\n')
	}
}

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
		ops := 0
		for i := 1; i+1 < n; i++ {
			if a[i] > a[i-1] && a[i] > a[i+1] {
				ops++
				if i+2 < n {
					if a[i] > a[i+2] {
						a[i+1] = a[i]
					} else {
						a[i+1] = a[i+2]
					}
				} else {
					a[i+1] = a[i]
				}
			}
		}
		fmt.Fprintln(writer, ops)
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, a[i])
		}
		writer.WriteByte('\n')
	}
}

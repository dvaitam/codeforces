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
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		d := make([]int, n)
		for i := n; i >= 1; i-- {
			pos := 0
			for j := 0; j < i; j++ {
				if a[j] == i {
					pos = j + 1
					break
				}
			}
			di := pos % i
			d[i-1] = di
			if di > 0 {
				prefix := append([]int{}, a[di:i]...)
				prefix = append(prefix, a[:di]...)
				copy(a[:i], prefix)
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, d[i])
		}
		writer.WriteByte('\n')
	}
}

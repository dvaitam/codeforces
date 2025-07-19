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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		start, end := 1, n
		for i := 0; i < k; i++ {
			for j := i; j < n; j += k {
				if i&1 == 1 {
					a[j] = end
					end--
				} else {
					a[j] = start
					start++
				}
			}
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

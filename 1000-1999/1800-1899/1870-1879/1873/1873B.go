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
		maxProd := int64(0)
		for i := 0; i < n; i++ {
			prod := int64(1)
			for j := 0; j < n; j++ {
				val := a[j]
				if j == i {
					val++
				}
				prod *= int64(val)
			}
			if prod > maxProd {
				maxProd = prod
			}
		}
		fmt.Fprintln(writer, maxProd)
	}
}

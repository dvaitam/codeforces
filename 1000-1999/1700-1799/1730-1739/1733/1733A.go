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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		sum := 0
		for r := 0; r < k; r++ {
			maxVal := 0
			for i := r; i < n; i += k {
				if a[i] > maxVal {
					maxVal = a[i]
				}
			}
			sum += maxVal
		}
		fmt.Fprintln(writer, sum)
	}
}

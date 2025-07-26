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
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		diff := make([]int64, n)
		maxDiff := int64(-1 << 63)
		for i := 0; i < n; i++ {
			diff[i] = a[i] - b[i]
			if diff[i] > maxDiff {
				maxDiff = diff[i]
			}
		}
		res := make([]int, 0)
		for i := 0; i < n; i++ {
			if diff[i] == maxDiff {
				res = append(res, i+1)
			}
		}
		fmt.Fprintln(writer, len(res))
		for i, v := range res {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
}

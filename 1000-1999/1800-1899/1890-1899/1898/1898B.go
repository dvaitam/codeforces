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
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var ops int64
		cur := a[n-1]
		for i := n - 2; i >= 0; i-- {
			if a[i] <= cur {
				cur = a[i]
			} else {
				k := (a[i] + cur - 1) / cur
				ops += k - 1
				cur = a[i] / k
			}
		}
		fmt.Fprintln(writer, ops)
	}
}

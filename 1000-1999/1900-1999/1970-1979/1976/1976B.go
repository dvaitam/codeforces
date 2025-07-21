package main

import (
	"bufio"
	"fmt"
	"math"
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
		b := make([]int, n+1)
		for i := 0; i <= n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		var base int64
		for i := 0; i < n; i++ {
			diff := int64(a[i] - b[i])
			if diff < 0 {
				diff = -diff
			}
			base += diff
		}
		best := int64(math.MaxInt64)
		for i := 0; i < n; i++ {
			diff := int64(a[i] - b[n])
			if diff < 0 {
				diff = -diff
			}
			if diff < best {
				best = diff
			}
		}
		fmt.Fprintln(writer, base+best)
	}
}

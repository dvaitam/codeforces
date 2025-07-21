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
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}
		best := int(^uint(0) >> 1)
		for i := 0; i < n-1; i++ {
			mx := a[i]
			if a[i+1] > mx {
				mx = a[i+1]
			}
			if mx < best {
				best = mx
			}
		}
		fmt.Fprintln(writer, best-1)
	}
}

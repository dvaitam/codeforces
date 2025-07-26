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
		cur := a[0]
		best := cur
		for i := 1; i < n; i++ {
			if (a[i]^a[i-1])&1 == 1 {
				if cur+a[i] > a[i] {
					cur += a[i]
				} else {
					cur = a[i]
				}
			} else {
				cur = a[i]
			}
			if cur > best {
				best = cur
			}
		}
		fmt.Fprintln(writer, best)
	}
}

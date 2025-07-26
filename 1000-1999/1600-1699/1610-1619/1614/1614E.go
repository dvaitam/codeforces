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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	T := make([]int64, n)
	lastans := int64(0)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &T[i])
		var k int
		fmt.Fscan(reader, &k)
		xs := make([]int64, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(reader, &xs[j])
		}
		for j := 0; j < k; j++ {
			x := (xs[j] + lastans) % (1000000000 + 1)
			cur := x
			for d := 0; d <= i; d++ {
				if cur < T[d] {
					cur++
				} else if cur > T[d] {
					cur--
				}
			}
			fmt.Fprintln(writer, cur)
			lastans = cur
		}
	}
}

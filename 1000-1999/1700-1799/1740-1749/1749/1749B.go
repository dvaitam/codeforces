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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int64, n)
		var sumB, maxB int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
			sumB += b[i]
			if b[i] > maxB {
				maxB = b[i]
			}
		}
		var sumA int64
		for _, v := range a {
			sumA += v
		}
		ans := sumA + sumB - maxB
		fmt.Fprintln(writer, ans)
	}
}

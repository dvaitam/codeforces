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
		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}

		dpPrev2 := 0
		dpPrev1 := freq[1]
		for i := 2; i <= n; i++ {
			current := dpPrev1
			if dpPrev2+freq[i] > current {
				current = dpPrev2 + freq[i]
			}
			dpPrev2 = dpPrev1
			dpPrev1 = current
		}
		keep := dpPrev1
		fmt.Fprintln(writer, n-keep)
	}
}

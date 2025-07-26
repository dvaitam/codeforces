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
		var cnt int
		for i := 0; i < n; i++ {
			for a[i]%2 == 0 {
				a[i] /= 2
				cnt++
			}
		}
		idx := 0
		for i := 1; i < n; i++ {
			if a[i] > a[idx] {
				idx = i
			}
		}
		for i := 0; i < cnt; i++ {
			a[idx] *= 2
		}
		var sum int64
		for i := 0; i < n; i++ {
			sum += a[i]
		}
		fmt.Fprintln(writer, sum)
	}
}

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
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		c := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &c[i])
		}
		count := 0
		for i := 0; i < n; i++ {
			bi := b[i]
			for j := 0; j < m; j++ {
				if bi+c[j] <= k {
					count++
				}
			}
		}
		fmt.Fprintln(writer, count)
	}
}

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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		banned := make([]bool, n+1)
		for i := 0; i < m; i++ {
			var a, b, c int
			fmt.Fscan(reader, &a, &b, &c)
			banned[b] = true
		}
		root := 1
		for root <= n && banned[root] {
			root++
		}
		for i := 1; i <= n; i++ {
			if i == root {
				continue
			}
			fmt.Fprintf(writer, "%d %d\n", root, i)
		}
	}
}

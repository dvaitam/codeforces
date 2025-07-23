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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	x := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x[i])
	}
	mark := make([]bool, 10)
	for i := 0; i < m; i++ {
		var y int
		fmt.Fscan(reader, &y)
		if y >= 0 && y < 10 {
			mark[y] = true
		}
	}

	first := true
	for i := 0; i < n; i++ {
		if mark[x[i]] {
			if !first {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, x[i])
			first = false
		}
	}
	writer.WriteByte('\n')
}

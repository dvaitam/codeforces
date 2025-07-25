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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var a, b string
		fmt.Fscan(reader, &a)
		fmt.Fscan(reader, &b)
		i, j := 0, 0
		for i < n && j < m {
			if a[i] == b[j] {
				i++
			}
			j++
		}
		fmt.Fprintln(writer, i)
	}
}

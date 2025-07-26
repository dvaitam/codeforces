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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		if n == 1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		q := make([]int, n)
		for i := 0; i < n; i++ {
			q[i] = i + 1
		}
		for i := 0; i < n; i++ {
			if q[i] == p[i] {
				if i+1 < n {
					q[i], q[i+1] = q[i+1], q[i]
				} else {
					q[i], q[i-1] = q[i-1], q[i]
				}
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, q[i])
		}
		writer.WriteByte('\n')
	}
}

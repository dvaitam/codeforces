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
	for t > 0 {
		t--
		var n, m int
		fmt.Fscan(reader, &n, &m)
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				a := (i%4 <= 1)
				b := (j%4 <= 1)
				if (a && b) || (!a && !b) {
					writer.WriteByte('1')
				} else {
					writer.WriteByte('0')
				}
				if j < m {
					writer.WriteByte(' ')
				}
			}
			writer.WriteByte('\n')
		}
	}
}

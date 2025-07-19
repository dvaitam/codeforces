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
		var n int
		fmt.Fscan(reader, &n)
		curr := 1
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(reader, &k)
			var out int
			if curr != 0 && (k%curr == 0 || curr%k == 0) {
				out = k
			} else {
				x := (k / curr) * curr
				if x < 1 {
					x = 1
				}
				out = x
			}
			curr = out
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, out)
		}
		writer.WriteByte('\n')
	}
}

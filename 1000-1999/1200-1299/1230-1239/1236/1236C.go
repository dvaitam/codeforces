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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var val int
			if j*2 < n {
				val = j*n + i + 1
			} else {
				val = j*n + n - i
			}
			if j > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, val)
		}
		writer.WriteByte('\n')
	}
}

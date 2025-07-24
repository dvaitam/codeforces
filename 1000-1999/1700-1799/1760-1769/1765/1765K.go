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
	fmt.Fscan(reader, &n)

	var total int64
	const inf int64 = 1<<63 - 1
	minDiag := inf

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var x int64
			fmt.Fscan(reader, &x)
			total += x
			if j == n-1-i && x < minDiag {
				minDiag = x
			}
		}
	}

	fmt.Fprintln(writer, total-minDiag)
}

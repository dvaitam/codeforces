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
	maxTime := 0
	const limit = 1000000
	for i := 0; i < n; i++ {
		var pos int
		fmt.Fscan(reader, &pos)
		left := pos - 1
		right := limit - pos
		if right < left {
			left = right
		}
		if left > maxTime {
			maxTime = left
		}
	}
	fmt.Fprintln(writer, maxTime)
}

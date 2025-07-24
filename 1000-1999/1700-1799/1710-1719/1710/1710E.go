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
	minA := int64(1<<63 - 1)
	for i := 0; i < n; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		if v < minA {
			minA = v
		}
	}
	var maxB int64
	for j := 0; j < m; j++ {
		var v int64
		fmt.Fscan(reader, &v)
		if j == 0 || v > maxB {
			maxB = v
		}
	}
	fmt.Fprintln(writer, minA+maxB)
}

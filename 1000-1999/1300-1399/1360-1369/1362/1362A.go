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
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		if a == b {
			fmt.Fprintln(writer, 0)
			continue
		}
		var bigger, smaller int64
		if a > b {
			bigger = a
			smaller = b
		} else {
			bigger = b
			smaller = a
		}
		if bigger%smaller != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		ratio := bigger / smaller
		var k int
		for ratio > 1 && ratio%2 == 0 {
			ratio /= 2
			k++
		}
		if ratio != 1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		ops := (k + 2) / 3
		fmt.Fprintln(writer, ops)
	}
}

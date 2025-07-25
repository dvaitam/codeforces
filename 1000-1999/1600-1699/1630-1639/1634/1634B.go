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
		var x, y int64
		fmt.Fscan(reader, &n, &x, &y)
		parity := int64(0)
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(reader, &a)
			parity ^= a & 1
		}
		aliceParity := (x & 1) ^ parity
		if aliceParity == (y & 1) {
			fmt.Fprintln(writer, "Alice")
		} else {
			fmt.Fprintln(writer, "Bob")
		}
	}
}

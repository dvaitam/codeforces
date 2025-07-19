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
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		var sum int64
		var xo int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			sum += x
			xo ^= x
		}
		fmt.Fprintln(writer, 2)
		fmt.Fprintln(writer, xo, sum+xo)
	}
}

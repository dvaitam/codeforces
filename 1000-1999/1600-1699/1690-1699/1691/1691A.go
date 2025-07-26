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
		even, odd := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x%2 == 0 {
				even++
			} else {
				odd++
			}
		}
		if even < odd {
			fmt.Fprintln(writer, even)
		} else {
			fmt.Fprintln(writer, odd)
		}
	}
}

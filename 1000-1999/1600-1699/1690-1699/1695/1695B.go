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
		minVal := int(^uint(0) >> 1)
		minIdx := 0
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			if v < minVal {
				minVal = v
				minIdx = i + 1
			}
		}
		if n%2 == 1 {
			fmt.Fprintln(writer, "Mike")
		} else {
			if minIdx%2 == 0 {
				fmt.Fprintln(writer, "Mike")
			} else {
				fmt.Fprintln(writer, "Joe")
			}
		}
	}
}

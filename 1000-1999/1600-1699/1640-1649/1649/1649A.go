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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		firstZero, lastZero := -1, -1
		for i, v := range arr {
			if v == 0 {
				if firstZero == -1 {
					firstZero = i
				}
				lastZero = i
			}
		}
		if firstZero == -1 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, lastZero-firstZero+2)
		}
	}
}

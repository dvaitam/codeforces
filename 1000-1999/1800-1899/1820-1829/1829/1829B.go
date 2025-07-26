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
		maxZero := 0
		curZero := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == 0 {
				curZero++
			} else {
				if curZero > maxZero {
					maxZero = curZero
				}
				curZero = 0
			}
		}
		if curZero > maxZero {
			maxZero = curZero
		}
		fmt.Fprintln(writer, maxZero)
	}
}

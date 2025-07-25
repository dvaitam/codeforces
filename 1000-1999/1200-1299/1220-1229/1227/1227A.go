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
		var n int
		fmt.Fscan(reader, &n)
		var maxL, minR int
		maxL = 0
		minR = 1<<31 - 1
		for i := 0; i < n; i++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			if l > maxL {
				maxL = l
			}
			if r < minR {
				minR = r
			}
		}
		if maxL <= minR {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, maxL-minR)
		}
	}
}

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
		segs := 0
		inSeg := false
		for i := 0; i < n; i++ {
			if arr[i] != 0 {
				if !inSeg {
					segs++
					inSeg = true
				}
			} else {
				inSeg = false
			}
		}
		if segs > 2 {
			segs = 2
		}
		fmt.Fprintln(writer, segs)
	}
}

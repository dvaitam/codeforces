package main

import (
	"bufio"
	"fmt"
	"math/bits"
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
		var W int64
		fmt.Fscan(reader, &n, &W)
		counts := make([]int, 31)
		for i := 0; i < n; i++ {
			var w int
			fmt.Fscan(reader, &w)
			idx := bits.TrailingZeros(uint(w))
			counts[idx]++
		}
		remaining := n
		height := 0
		for remaining > 0 {
			height++
			remW := W
			for j := 30; j >= 0; j-- {
				if counts[j] == 0 {
					continue
				}
				width := int64(1 << j)
				if width > remW {
					continue
				}
				maxFit := int(remW / width)
				if maxFit > counts[j] {
					maxFit = counts[j]
				}
				counts[j] -= maxFit
				remW -= width * int64(maxFit)
				remaining -= maxFit
			}
		}
		fmt.Fprintln(writer, height)
	}
}

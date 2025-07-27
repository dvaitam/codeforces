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
		maxFreq := 0
		prev := -1
		freq := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x == prev {
				freq++
			} else {
				if freq > maxFreq {
					maxFreq = freq
				}
				prev = x
				freq = 1
			}
		}
		if freq > maxFreq {
			maxFreq = freq
		}
		fmt.Fprintln(writer, maxFreq)
	}
}

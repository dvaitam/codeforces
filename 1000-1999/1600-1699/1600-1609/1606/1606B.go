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
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		hours := int64(0)
		have := int64(1)
		for have < k && have < n {
			have *= 2
			hours++
		}
		if have < n {
			remaining := n - have
			hours += (remaining + k - 1) / k
		}
		fmt.Fprintln(writer, hours)
	}
}

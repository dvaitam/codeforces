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
		var r, g, b int64
		fmt.Fscan(reader, &r, &g, &b)
		sum := r + g + b
		max := r
		if g > max {
			max = g
		}
		if b > max {
			max = b
		}
		ans := sum - max
		if sum/2 < ans {
			ans = sum / 2
		}
		fmt.Fprintln(writer, ans)
	}
}

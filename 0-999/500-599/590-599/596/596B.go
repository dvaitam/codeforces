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

	var n int
	fmt.Fscan(reader, &n)
	var prev int64
	var ans int64
	for i := 0; i < n; i++ {
		var b int64
		fmt.Fscan(reader, &b)
		diff := b - prev
		if diff < 0 {
			diff = -diff
		}
		ans += diff
		prev = b
	}
	fmt.Fprintln(writer, ans)
}

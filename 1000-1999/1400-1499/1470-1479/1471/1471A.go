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
		var x int64
		fmt.Fscan(reader, &n, &x)
		var sum int64
		var maxBeauty int64
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(reader, &v)
			sum += v
			maxBeauty += (v + x - 1) / x
		}
		minBeauty := (sum + x - 1) / x
		fmt.Fprintf(writer, "%d %d\n", minBeauty, maxBeauty)
	}
}

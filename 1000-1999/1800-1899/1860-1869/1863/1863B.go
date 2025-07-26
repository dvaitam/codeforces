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
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			pos[x] = i
		}
		segments := 1
		for v := 2; v <= n; v++ {
			if pos[v] < pos[v-1] {
				segments++
			}
		}
		fmt.Fprintln(writer, segments-1)
	}
}

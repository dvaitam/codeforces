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
		ans := make([]int, n+1)
		used := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			if used[i] {
				continue
			}
			fmt.Fprintln(writer, "?", i)
			writer.Flush()

			var tmp int
			fmt.Fscan(reader, &tmp)

			fmt.Fprintln(writer, "?", i)
			writer.Flush()

			var cur int
			fmt.Fscan(reader, &cur)
			cycle := []int{cur}
			for {
				fmt.Fprintln(writer, "?", i)
				writer.Flush()
				fmt.Fscan(reader, &cur)
				if cur == cycle[0] {
					break
				}
				cycle = append(cycle, cur)
			}
			for j := 0; j < len(cycle); j++ {
				nxt := cycle[(j+1)%len(cycle)]
				ans[cycle[j]] = nxt
				used[cycle[j]] = true
			}
		}
		fmt.Fprint(writer, "!")
		for i := 1; i <= n; i++ {
			fmt.Fprint(writer, " ", ans[i])
		}
		fmt.Fprintln(writer)
		writer.Flush()
	}
}

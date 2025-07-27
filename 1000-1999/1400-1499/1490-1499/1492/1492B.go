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
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
			pos[p[i]] = i
		}

		last := n - 1
		val := n
		first := true
		for last >= 0 {
			for val > 0 && pos[val] > last {
				val--
			}
			start := pos[val]
			for i := start; i <= last; i++ {
				if !first {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, p[i])
				first = false
			}
			last = start - 1
			val--
		}
		fmt.Fprintln(writer)
	}
}

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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	weights := make(map[int64]int64)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v, u, w int64
			fmt.Fscan(reader, &v, &u, &w)
			for v != u {
				if v > u {
					weights[v] += w
					v /= 2
				} else {
					weights[u] += w
					u /= 2
				}
			}
		} else if t == 2 {
			var v, u int64
			fmt.Fscan(reader, &v, &u)
			var res int64
			for v != u {
				if v > u {
					res += weights[v]
					v /= 2
				} else {
					res += weights[u]
					u /= 2
				}
			}
			fmt.Fprintln(writer, res)
		}
	}
}

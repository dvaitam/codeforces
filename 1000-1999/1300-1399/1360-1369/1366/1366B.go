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
		var n, x, m int
		fmt.Fscan(reader, &n, &x, &m)
		l, r := x, x
		for i := 0; i < m; i++ {
			var li, ri int
			fmt.Fscan(reader, &li, &ri)
			if li <= r && ri >= l {
				if li < l {
					l = li
				}
				if ri > r {
					r = ri
				}
			}
		}
		fmt.Fprintln(writer, r-l+1)
	}
}

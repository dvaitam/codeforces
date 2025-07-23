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
	for ; q > 0; q-- {
		var n, m, p int
		fmt.Fscan(reader, &n, &m, &p)
		if p > n*m {
			fmt.Fprintln(writer, "No")
			continue
		}
		possible := false
		for d := 1; d*d <= p && !possible; d++ {
			if p%d == 0 {
				x := d
				y := p / d
				if (x <= n && y <= m) || (x <= m && y <= n) {
					possible = true
				}
			}
		}
		if possible {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}

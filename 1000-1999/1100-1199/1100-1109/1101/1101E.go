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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var mxx, mxy int64
	for i := 0; i < n; i++ {
		var op string
		var x, y int64
		fmt.Fscan(reader, &op, &x, &y)
		if x > y {
			x, y = y, x
		}
		if op[0] == '+' {
			if x > mxx {
				mxx = x
			}
			if y > mxy {
				mxy = y
			}
		} else {
			if mxx <= x && mxy <= y {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}

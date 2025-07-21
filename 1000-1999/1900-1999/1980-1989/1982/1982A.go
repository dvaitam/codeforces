package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x1, y1, x2, y2 int64
		fmt.Fscan(in, &x1, &y1)
		fmt.Fscan(in, &x2, &y2)
		diff1 := x1 - y1
		diff2 := x2 - y2
		if diff1*diff2 > 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

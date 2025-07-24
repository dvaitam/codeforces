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
		var n, m, sx, sy, d int
		fmt.Fscan(in, &n, &m, &sx, &sy, &d)
		path1 := (sx-1 > d) && (m-sy > d)
		path2 := (sy-1 > d) && (n-sx > d)
		if path1 || path2 {
			fmt.Fprintln(out, n+m-2)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}

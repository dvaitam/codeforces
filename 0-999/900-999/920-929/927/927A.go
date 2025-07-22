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

	var w, h, k int
	if _, err := fmt.Fscan(in, &w, &h); err != nil {
		return
	}
	fmt.Fscan(in, &k)
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
	}

	fmt.Fprintln(out, 0)
	out.Flush()

	for {
		var t int
		if _, err := fmt.Fscan(in, &t); err != nil {
			return
		}
		var sx, sy, tx, ty int
		fmt.Fscan(in, &sx, &sy, &tx, &ty)
		if t == -1 {
			fmt.Fprintln(out, 0)
			out.Flush()
			return
		}
		fmt.Fprintln(out, 0)
		out.Flush()
	}
}

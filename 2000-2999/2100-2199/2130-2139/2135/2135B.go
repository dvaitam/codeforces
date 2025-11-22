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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n; i++ {
			var x, y int64
			fmt.Fscan(in, &x, &y)
		}
		var X, Y int64
		fmt.Fscan(in, &X, &Y)
		fmt.Fprint(out, X, " ", Y)
		if T > 1 {
			fmt.Fprintln(out)
		}
	}
}

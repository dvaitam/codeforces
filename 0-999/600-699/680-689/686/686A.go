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

	var n int
	var x int64
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	distressed := 0
	for i := 0; i < n; i++ {
		var op string
		var d int64
		fmt.Fscan(in, &op, &d)
		if op == "+" {
			x += d
		} else {
			if x >= d {
				x -= d
			} else {
				distressed++
			}
		}
	}
	fmt.Fprintf(out, "%d %d", x, distressed)
}

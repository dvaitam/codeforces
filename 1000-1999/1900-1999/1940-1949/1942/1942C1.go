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
		var n, x, y int
		fmt.Fscan(in, &n, &x, &y)
		for i := 0; i < x; i++ {
			var tmp int
			fmt.Fscan(in, &tmp)
		}
		if x < 3 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, x-2)
		}
	}
}

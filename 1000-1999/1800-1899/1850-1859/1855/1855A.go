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
		var n int
		fmt.Fscan(in, &n)
		fixed := 0
		for i := 1; i <= n; i++ {
			var p int
			fmt.Fscan(in, &p)
			if p == i {
				fixed++
			}
		}
		moves := fixed / 2
		if fixed%2 == 1 {
			moves++
		}
		fmt.Fprintln(out, moves)
	}
}

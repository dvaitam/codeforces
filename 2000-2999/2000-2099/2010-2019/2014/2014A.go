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
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		gold := 0
		given := 0
		for i := 0; i < n; i++ {
			var ai int
			fmt.Fscan(in, &ai)
			if ai >= k {
				gold += ai
			} else if ai == 0 {
				if gold > 0 {
					gold--
					given++
				}
			}
		}
		fmt.Fprintln(out, given)
	}
}

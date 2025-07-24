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
		pos := 1
		for i := 1; i <= n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v == 0 {
				pos = i
			}
		}
		fmt.Fprintln(out, pos, pos)
	}
}

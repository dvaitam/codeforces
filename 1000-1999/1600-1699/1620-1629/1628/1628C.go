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
		res := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				var x int
				fmt.Fscan(in, &x)
				res ^= x
			}
		}
		fmt.Fprintln(out, res)
	}
}

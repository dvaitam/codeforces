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
		seen := make([]bool, 101)
		distinct := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if !seen[x] {
				seen[x] = true
				distinct++
			}
		}
		answer := 2*distinct - 1
		fmt.Fprintln(out, answer)
	}
}

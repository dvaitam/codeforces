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
		var n, x int
		fmt.Fscan(in, &n, &x)
		present := make(map[int]bool, n)
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			present[v] = true
		}
		for i := 1; ; i++ {
			if !present[i] {
				if x > 0 {
					x--
				} else {
					fmt.Fprintln(out, i-1)
					break
				}
			}
		}
	}
}

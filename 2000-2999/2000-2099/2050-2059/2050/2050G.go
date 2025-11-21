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
		var n int
		fmt.Fscan(in, &n)
		maxDeg := 0
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			deg := make([]int, n+1)
			for j := 0; j < n-1; j++ {
				deg[u]++
				deg[v]++
				if deg[u] > maxDeg {
					maxDeg = deg[u]
				}
				if deg[v] > maxDeg {
					maxDeg = deg[v]
				}
			}
		}
		fmt.Fprintln(out, maxDeg)
	}
}

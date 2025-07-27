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
		var n, x int
		fmt.Fscan(in, &n, &x)
		deg := make([]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			deg[u]++
			deg[v]++
		}
		if deg[x] <= 1 || n%2 == 0 {
			fmt.Fprintln(out, "Ayush")
		} else {
			fmt.Fprintln(out, "Ashish")
		}
	}
}

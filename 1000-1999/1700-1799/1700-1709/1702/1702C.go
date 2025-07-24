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
		var n, k int
		fmt.Fscan(in, &n, &k)
		first := make(map[int]int)
		last := make(map[int]int)
		for i := 0; i < n; i++ {
			var u int
			fmt.Fscan(in, &u)
			if _, ok := first[u]; !ok {
				first[u] = i
			}
			last[u] = i
		}
		for ; k > 0; k-- {
			var a, b int
			fmt.Fscan(in, &a, &b)
			fa, okA := first[a]
			lb, okB := last[b]
			if okA && okB && fa <= lb {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

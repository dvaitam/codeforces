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
		type edge struct{ u, v int }
		edges := make([]edge, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v)
		}

		fmt.Fprint(out, "!")
		for _, e := range edges {
			fmt.Fprintf(out, " %d %d", e.u, e.v)
		}
		fmt.Fprintln(out)
	}
}

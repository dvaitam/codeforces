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
		parents := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &parents[i])
		}
		fmt.Fprint(out, "!")
		for i := 0; i < n-1; i++ {
			fmt.Fprintf(out, " %d", parents[i])
		}
		fmt.Fprintln(out)
	}
}

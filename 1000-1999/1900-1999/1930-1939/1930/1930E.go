package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the actual algorithm for problem E.
// This placeholder reads the input and outputs zero for each k.

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
		m := (n - 1) / 2
		for i := 0; i < m; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, 0)
		}
		out.WriteByte('\n')
	}
}

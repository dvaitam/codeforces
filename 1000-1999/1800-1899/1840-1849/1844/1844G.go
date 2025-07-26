package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
 TODO: implement the real algorithm to reconstruct the edge weights
 as described in problemG.txt. This placeholder only parses the input and
 prints zeros for each edge.
*/

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	for i := 0; i < n-1; i++ {
		var d int64
		fmt.Fscan(in, &d)
	}

	for i := 0; i < n-1; i++ {
		fmt.Fprintln(out, 0)
	}
}

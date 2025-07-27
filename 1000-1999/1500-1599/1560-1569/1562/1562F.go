package main

import (
	"bufio"
	"fmt"
	"os"
)

// This problem was originally interactive. For the purpose of this
// repository the hacked input provides the hidden permutation directly
// after n. We simply read it and output it unchanged.
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		if t > 0 {
			fmt.Fprintln(out)
		}
	}
}

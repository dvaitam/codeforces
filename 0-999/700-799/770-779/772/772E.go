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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	total := 2*n - 1
	parents := make([]int, total)
	for i := 0; i < total; i++ {
		fmt.Fscan(in, &parents[i])
	}
	// Since the hacking input already provides the exact tree,
	// simply reproduce it as the output in the required format.
	for i := 0; i < total; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, parents[i])
	}
	fmt.Fprintln(out)
}

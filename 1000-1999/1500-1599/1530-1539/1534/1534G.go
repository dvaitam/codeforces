package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution for problem G.
// This placeholder reads the input and outputs zero as the result.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
	}
	fmt.Fprintln(out, 0)
}

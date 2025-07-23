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
	weights := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
		total += weights[i]
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		// edges are ignored in this simple solution
		_ = u
		_ = v
	}
	fmt.Fprintln(out, total)
}

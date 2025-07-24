package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for Codeforces 1718F.
// TODO: implement the actual algorithm.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, C, q int
	if _, err := fmt.Fscan(in, &n, &m, &C, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		// The computation is not implemented yet.
		fmt.Fprintln(out, 0)
	}
}

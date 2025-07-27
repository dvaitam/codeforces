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

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	w := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &w[i])
	}
	for i := 0; i < m; i++ {
		var u, v int
		var d int64
		fmt.Fscan(in, &u, &v, &d)
	}
	// initial cost using station 1 for all connections
	base := int64(n-1) * w[1]
	fmt.Fprintln(out, base)
	for i := 0; i < q; i++ {
		var k int
		var x int64
		fmt.Fscan(in, &k, &x)
		w[k] += x
		if k == 1 {
			base = int64(n-1) * w[1]
		}
		fmt.Fprintln(out, base)
	}
}

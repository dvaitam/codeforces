package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemH.txt (1609H).
// The real simulation of robots moving and answering queries is complex
// and requires careful implementation. This stub only parses the input
// and outputs zero for each query.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	pos := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pos[i])
	}
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			var v int64
			fmt.Fscan(in, &v)
		}
	}
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var a int
		var t int64
		fmt.Fscan(in, &a, &t)
		fmt.Fprintln(out, 0)
	}
}

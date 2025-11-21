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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		p := make([]int, n)
		for i := 2; i <= n; i++ {
			fmt.Fscan(in, &p[i-1])
		}
		fmt.Fprintln(out, minInt(k, n-k))
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

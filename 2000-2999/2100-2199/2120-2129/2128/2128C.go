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
		var n int
		fmt.Fscan(in, &n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		ok := true
		prefixMin := b[0]
		for i := 1; i < n && ok; i++ {
			limit := prefixMin * (prefixMin + 1) / 2
			if b[i] > limit {
				ok = false
				break
			}
			if b[i] < prefixMin {
				prefixMin = b[i]
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		ans := n
		for m := 1; m <= n; m++ {
			if m*(n-m) <= k {
				ans = m
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}

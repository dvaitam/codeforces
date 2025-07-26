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
		var n, k int64
		fmt.Fscan(in, &n, &k)
		var ans int64
		if k >= 31 {
			ans = n + 1
		} else {
			pow := int64(1) << k
			if n+1 < pow {
				ans = n + 1
			} else {
				ans = pow
			}
		}
		fmt.Fprintln(out, ans)
	}
}

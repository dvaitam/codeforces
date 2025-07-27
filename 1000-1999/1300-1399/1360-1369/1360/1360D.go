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
		if k >= n {
			fmt.Fprintln(out, 1)
			continue
		}
		ans := n
		for i := 1; i*i <= n; i++ {
			if n%i == 0 {
				d1 := i
				d2 := n / i
				if d1 <= k {
					if n/d1 < ans {
						ans = n / d1
					}
				}
				if d2 <= k {
					if n/d2 < ans {
						ans = n / d2
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

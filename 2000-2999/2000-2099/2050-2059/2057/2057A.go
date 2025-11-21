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
		var n, m int64
		fmt.Fscan(in, &n, &m)

		var ans int64
		if n < m {
			n, m = m, n
		}

		if m == 1 {
			ans = n*(n-1)/2 + n + n
		} else {
			ans = m*(m-1)/2 + m + (m+1)*(m+2)/2
			if n > m {
				ans += (n - (m + 1)) * (m + 1)
			}
		}
		fmt.Fprintln(out, ans)
	}
}

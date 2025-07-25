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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int64
		fmt.Fscan(in, &n, &k)

		if k > n {
			fmt.Fprintln(out, "NO")
			continue
		}
		if k == n {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, 1)
			continue
		}
		if k <= (n+1)/2 {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, 2)
			fmt.Fprintln(out, n-k+1, 1)
			continue
		}
		fmt.Fprintln(out, "NO")
	}
}

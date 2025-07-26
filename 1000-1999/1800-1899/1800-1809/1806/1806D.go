package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// TODO: implement a correct solution. This placeholder outputs zeros.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		for i := 0; i < n-1; i++ {
			var tmp int
			fmt.Fscan(in, &tmp)
		}
		for k := 1; k < n; k++ {
			if k > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, 0)
		}
		fmt.Fprintln(out)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 1000000

func main() {
	cntDiv := make([]int, MAXN+1)
	for d := 1; d <= MAXN; d++ {
		for m := d; m <= MAXN; m += d {
			cntDiv[m]++
		}
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		b := 0
		for i := 0; i < n; i++ {
			if i > 0 && s[i] != s[i-1] {
				b++
			}
			var ans int
			if b == 0 {
				ans = i + 1
			} else {
				ans = cntDiv[b]
			}
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}

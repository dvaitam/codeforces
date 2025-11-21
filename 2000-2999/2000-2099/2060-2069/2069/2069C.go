package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var ans, sumG, cntOne int64
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			switch x {
			case 1:
				cntOne++
				if cntOne >= mod {
					cntOne -= mod
				}
			case 2:
				sumG = (2*sumG + cntOne) % mod
			case 3:
				ans += sumG
				if ans >= mod {
					ans -= mod
				}
			default:
				// values limited to 1..3
			}
		}
		fmt.Fprintln(out, ans%mod)
	}
}

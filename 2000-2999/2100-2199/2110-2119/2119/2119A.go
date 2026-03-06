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
		var a, b int
		var x, y int64
		fmt.Fscan(in, &a, &b, &x, &y)

		if a == b {
			fmt.Fprintln(out, 0)
			continue
		}

		if a > b {
			// Can only decrease by 1 via XOR when a is odd (odd XOR 1 = odd-1)
			if a%2 == 1 && b == a-1 {
				fmt.Fprintln(out, y)
			} else {
				fmt.Fprintln(out, -1)
			}
			continue
		}

		// a < b: move from a up to b
		// At even k: both +1 and XOR go to k+1, so cost is min(x, y)
		// At odd k: XOR goes to k-1 (wrong direction), so must use +1, cost x
		var ans int64
		for k := a; k < b; k++ {
			if k%2 == 0 {
				if x < y {
					ans += x
				} else {
					ans += y
				}
			} else {
				ans += x
			}
		}
		fmt.Fprintln(out, ans)
	}
}

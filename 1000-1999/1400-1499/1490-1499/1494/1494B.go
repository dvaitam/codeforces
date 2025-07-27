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
		var n, u, r, d, l int
		fmt.Fscan(in, &n, &u, &r, &d, &l)
		ok := false
		for mask := 0; mask < 16; mask++ {
			uu, rr, dd, ll := u, r, d, l
			if mask&1 != 0 {
				uu--
				ll--
			}
			if mask&2 != 0 {
				uu--
				rr--
			}
			if mask&4 != 0 {
				dd--
				rr--
			}
			if mask&8 != 0 {
				dd--
				ll--
			}
			if uu >= 0 && uu <= n-2 && rr >= 0 && rr <= n-2 && dd >= 0 && dd <= n-2 && ll >= 0 && ll <= n-2 {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		pref := 0
		var res int64
		for r := 1; r <= n; r++ {
			val := r - a[r] + 1
			if val > pref {
				pref = val
			}
			L := pref
			if L < 1 {
				L = 1
			}
			if L <= r {
				res += int64(r - L + 1)
			}
		}
		fmt.Fprintln(out, res)
	}
}

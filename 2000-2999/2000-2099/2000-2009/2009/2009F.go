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
		var n, q int
		fmt.Fscan(in, &n, &q)

		a := make([]int64, n)
		total := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
		}

		pref := make([]int64, 2*n+1)
		for i := 1; i <= 2*n; i++ {
			pref[i] = pref[i-1] + a[(i-1)%n]
		}

		part := func(start, length int) int64 {
			if length <= 0 {
				return 0
			}
			l := start
			r := start + length - 1
			return pref[r] - pref[l-1]
		}

		var prefix func(int64) int64
		prefix = func(pos int64) int64 {
			if pos <= 0 {
				return 0
			}
			full := int(pos / int64(n))
			rem := int(pos % int64(n))
			sum := int64(full) * total
			if rem == 0 {
				return sum
			}
			shiftIdx := full + 1
			if shiftIdx > n {
				shiftIdx = n
			}
			sum += part(shiftIdx, rem)
			return sum
		}

		for ; q > 0; q-- {
			var l, r int64
			fmt.Fscan(in, &l, &r)
			ans := prefix(r) - prefix(l-1)
			fmt.Fprintln(out, ans)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Solution for problemB.txt of contest 1497 (M-arrays).
// We group numbers by their remainder modulo m. For each
// remainder r we consider its complement m-r. Numbers with
// r and m-r can be arranged in alternating order so that
// adjacent sums are divisible by m. The minimal number of
// such arrays is obtained by combining each pair (r,m-r)
// greedily: if their counts differ by at most one they form
// a single array, otherwise after forming one alternating
// array we still have diff-1 leftover elements from the
// larger remainder which must be single-element arrays.
// Remainders 0 and m/2 (when m is even) are handled
// separately because they pair with themselves.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		cnt := make([]int, m)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x%m]++
		}
		ans := 0
		if cnt[0] > 0 {
			ans++
		}
		for r := 1; r*2 < m; r++ {
			a := cnt[r]
			b := cnt[m-r]
			if a == 0 && b == 0 {
				continue
			}
			diff := int(math.Abs(float64(a - b)))
			if diff <= 1 {
				ans++
			} else {
				ans += diff
			}
		}
		if m%2 == 0 {
			if cnt[m/2] > 0 {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}

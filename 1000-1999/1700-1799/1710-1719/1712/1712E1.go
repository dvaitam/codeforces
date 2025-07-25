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
		var l, r int
		fmt.Fscan(in, &l, &r)
		n := int64(r - l + 1)
		total := n * (n - 1) * (n - 2) / 6

		divCount := make([]int, r+1)
		for d := l; d <= r; d++ {
			start := (l + d - 1) / d * d
			if start <= d {
				start += d
			}
			for k := start; k <= r; k += d {
				divCount[k]++
			}
		}

		var fail1 int64
		for k := l + 2; k <= r; k++ {
			m := int64(divCount[k])
			fail1 += m * (m - 1) / 2
		}

		var fail2 int64
		for k := l + 2; k <= r; k++ {
			if k%6 == 0 {
				i := k / 2
				j := 2 * k / 3
				if i >= l && j >= l && i < j && j < k {
					fail2++
				}
			}
			if k%15 == 0 {
				i := 2 * k / 5
				j := 2 * k / 3
				if i >= l && j >= l && i < j && j < k {
					fail2++
				}
			}
		}

		ans := total - fail1 - fail2
		fmt.Fprintln(out, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func ceilDiv(a, b int64) int64 {
	return (a + b - 1) / b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		var sum int64
		var mx int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
			if a[i] > mx {
				mx = a[i]
			}
		}
		best := 1
		for s := 1; s <= n; s++ {
			s64 := int64(s)
			d := mx
			need := ceilDiv(sum, s64)
			if need > d {
				d = need
			}
			add := d*s64 - sum
			if add <= k {
				if s > best {
					best = s
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}

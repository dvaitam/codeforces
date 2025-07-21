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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x, y, k int64
		fmt.Fscan(in, &x, &y, &k)
		for k > 0 && x >= y {
			rem := int64(y-1) - x%y
			if rem >= k {
				x += k
				k = 0
				break
			}
			x += rem
			k -= rem
			x++
			k--
			for x%y == 0 {
				x /= y
			}
		}
		if k > 0 {
			period := int64(y - 1)
			if period > 0 {
				k %= period
				x = ((x - 1 + k) % period) + 1
			}
		}
		fmt.Fprintln(out, x)
	}
}

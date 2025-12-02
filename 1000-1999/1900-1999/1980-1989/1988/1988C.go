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

	for i := 0; i < t; i++ {
		var n int64
		fmt.Fscan(in, &n)

		var bits []int64
		temp := n
		var p int64 = 1
		for temp > 0 {
			if temp&1 == 1 {
				bits = append(bits, p)
			}
			temp >>= 1
			p <<= 1
		}
		
		var ans []int64
		// Iterate bits from largest to smallest
		for j := len(bits) - 1; j >= 0; j-- {
			val := n - bits[j]
			if val > 0 {
				ans = append(ans, val)
			}
		}
		ans = append(ans, n)

		fmt.Fprintln(out, len(ans))
		for j, val := range ans {
			if j > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val)
		}
		fmt.Fprintln(out)
	}
}
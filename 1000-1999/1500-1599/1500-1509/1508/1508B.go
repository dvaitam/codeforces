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
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		limit := int64(1)
		for i := 0; i < n-1 && limit <= 1e18; i++ {
			limit <<= 1
		}
		if limit < k {
			fmt.Fprintln(out, -1)
			continue
		}

		bits := make([]bool, n-1)
		km1 := k - 1
		for i := 0; i < n-1; i++ {
			shift := n - 2 - i
			var b int64
			if shift >= 0 && shift < 60 {
				b = (km1 >> uint(shift)) & 1
			} else if shift >= 60 {
				b = 0
			}
			bits[i] = b == 0
		}

		first := true
		for i := 1; i <= n; {
			j := i
			for j <= n-1 && !bits[j-1] {
				j++
			}
			for x := j; x >= i; x-- {
				if !first {
					fmt.Fprint(out, " ")
				}
				first = false
				fmt.Fprint(out, x)
			}
			i = j + 1
		}
		fmt.Fprintln(out)
	}
}

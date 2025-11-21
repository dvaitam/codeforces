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
		var m int64
		fmt.Fscan(in, &n, &m)

		var prevTime int64
		var prevSide int64
		var ans int64

		for i := 0; i < n; i++ {
			var a, b int64
			fmt.Fscan(in, &a, &b)

			length := a - prevTime
			parity := prevSide ^ b
			if (length & 1) == parity {
				ans += length
			} else {
				ans += length - 1
			}

			prevTime = a
			prevSide = b
		}

		ans += m - prevTime
		fmt.Fprintln(out, ans)
	}
}

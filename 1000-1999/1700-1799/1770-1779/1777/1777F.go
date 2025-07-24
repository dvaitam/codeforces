package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit = 130

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefix := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] ^ a[i]
		}

		ans := 0
		for l := 0; l < n; l++ {
			mx := 0
			for r := l; r < n && r-l < limit; r++ {
				if a[r] > mx {
					mx = a[r]
				}
				x := prefix[r+1] ^ prefix[l]
				val := mx ^ x
				if val > ans {
					ans = val
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

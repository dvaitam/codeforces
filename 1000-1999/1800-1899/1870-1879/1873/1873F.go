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
		a := make([]int64, n)
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}

		ans := 0
		l := 0
		var sum int64
		for r := 0; r < n; r++ {
			if r > 0 && h[r-1]%h[r] != 0 {
				l = r
				sum = 0
			}
			sum += a[r]
			for sum > k {
				sum -= a[l]
				l++
			}
			if r-l+1 > ans {
				ans = r - l + 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}

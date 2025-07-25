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
		c := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &c[i])
		}
		prefixSum := make([]int64, n+1)
		prefixOnes := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefixSum[i] = prefixSum[i-1] + c[i]
			prefixOnes[i] = prefixOnes[i-1]
			if c[i] == 1 {
				prefixOnes[i]++
			}
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			length := r - l + 1
			sum := prefixSum[r] - prefixSum[l-1]
			ones := prefixOnes[r] - prefixOnes[l-1]
			if length <= 1 || sum < int64(length+ones) {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		}
	}
}

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
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		total := 0
		sumSquares := 0
		for i := 0; i < n; i++ {
			total += a[i] + b[i]
			sumSquares += a[i]*a[i] + b[i]*b[i]
		}
		dp := make([]bool, total+1)
		dp[0] = true
		for i := 0; i < n; i++ {
			ndp := make([]bool, total+1)
			for s := 0; s <= total; s++ {
				if dp[s] {
					if s+a[i] <= total {
						ndp[s+a[i]] = true
					}
					if s+b[i] <= total {
						ndp[s+b[i]] = true
					}
				}
			}
			dp = ndp
		}
		best := int64(1<<63 - 1)
		for s := 0; s <= total; s++ {
			if dp[s] {
				sumA := int64(s)
				sumB := int64(total - s)
				cur := sumA*sumA + sumB*sumB
				if cur < best {
					best = cur
				}
			}
		}
		constPart := int64(sumSquares) * int64(n-2)
		fmt.Fprintln(out, best+constPart)
	}
}

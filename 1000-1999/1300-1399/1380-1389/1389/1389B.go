package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k, z int
		fmt.Fscan(reader, &n, &k, &z)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sum := 0
		bestPair := 0
		ans := 0
		for i := 0; i <= k; i++ {
			sum += a[i]
			if i < n-1 {
				if v := a[i] + a[i+1]; v > bestPair {
					bestPair = v
				}
			}
			remaining := k - i
			left := z
			if remaining/2 < left {
				left = remaining / 2
			}
			candidate := sum + bestPair*left
			if candidate > ans {
				ans = candidate
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var q int
	fmt.Fscan(reader, &q)

	limit := int(math.Sqrt(float64(n))) + 1
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, limit+1)
	}

	for k := 1; k <= limit; k++ {
		for i := n; i >= 1; i-- {
			next := i + a[i] + k
			if next > n {
				dp[i][k] = 1
			} else {
				dp[i][k] = dp[next][k] + 1
			}
		}
	}

	for ; q > 0; q-- {
		var p, k int
		fmt.Fscan(reader, &p, &k)
		if k <= limit {
			fmt.Fprintln(writer, dp[p][k])
		} else {
			cnt := 0
			for p <= n {
				p = p + a[p] + k
				cnt++
			}
			fmt.Fprintln(writer, cnt)
		}
	}
}

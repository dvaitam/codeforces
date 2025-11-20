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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	type item struct {
		w int
		v int64
	}
	items := make([]item, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &items[i].w, &items[i].v)
	}

	const negInf = int64(-1 << 60)
	dp := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		dp[i] = negInf
	}
	dp[0] = 0

	for _, it := range items {
		w := it.w
		for weight := m; weight >= w; weight-- {
			if dp[weight-w] == negInf {
				continue
			}
			val := dp[weight-w] + it.v
			if val > dp[weight] {
				dp[weight] = val
			}
		}
	}

	ans := int64(0)
	for w := 0; w <= m; w++ {
		if dp[w] > ans {
			ans = dp[w]
		}
	}
	fmt.Fprintln(writer, ans)
}

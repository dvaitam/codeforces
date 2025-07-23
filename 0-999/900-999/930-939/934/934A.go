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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}

	best := make([]int64, n)
	for i := 0; i < n; i++ {
		maxVal := int64(math.MinInt64)
		for j := 0; j < m; j++ {
			prod := a[i] * b[j]
			if prod > maxVal {
				maxVal = prod
			}
		}
		best[i] = maxVal
	}

	ans := int64(math.MaxInt64)
	for hide := 0; hide < n; hide++ {
		maxVal := int64(math.MinInt64)
		for i := 0; i < n; i++ {
			if i == hide {
				continue
			}
			if best[i] > maxVal {
				maxVal = best[i]
			}
		}
		if maxVal < ans {
			ans = maxVal
		}
	}
	fmt.Fprintln(writer, ans)
}

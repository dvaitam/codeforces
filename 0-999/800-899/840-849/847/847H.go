package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	incVal := make([]int64, n+2)
	incCost := make([]int64, n+2)
	incVal[0] = -INF
	for i := 1; i <= n; i++ {
		if incVal[i-1]+1 > a[i] {
			incVal[i] = incVal[i-1] + 1
		} else {
			incVal[i] = a[i]
		}
		incCost[i] = incCost[i-1] + incVal[i] - a[i]
	}

	decVal := make([]int64, n+2)
	decCost := make([]int64, n+2)
	decVal[n+1] = -INF
	for i := n; i >= 1; i-- {
		if decVal[i+1]+1 > a[i] {
			decVal[i] = decVal[i+1] + 1
		} else {
			decVal[i] = a[i]
		}
		decCost[i] = decCost[i+1] + decVal[i] - a[i]
	}

	ans := INF
	for i := 1; i <= n; i++ {
		peak := incVal[i]
		if decVal[i] > peak {
			peak = decVal[i]
		}
		cost := incCost[i-1] + decCost[i+1] + peak - a[i]
		if cost < ans {
			ans = cost
		}
	}

	fmt.Println(ans)
}

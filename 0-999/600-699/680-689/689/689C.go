package main

import (
	"bufio"
	"fmt"
	"os"
)

func countWays(n int64) int64 {
	var res int64
	for k := int64(2); k*k*k <= n; k++ {
		res += n / (k * k * k)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var m int64
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}

	low := int64(1)
	high := int64(1e18)
	for low < high {
		mid := (low + high) / 2
		if countWays(mid) >= m {
			high = mid
		} else {
			low = mid + 1
		}
	}

	if countWays(low) == m {
		fmt.Println(low)
	} else {
		fmt.Println(-1)
	}
}

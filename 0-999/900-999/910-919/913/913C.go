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

	var n int
	var L int64
	if _, err := fmt.Fscan(reader, &n, &L); err != nil {
		return
	}
	const maxBits = 31
	costs := make([]int64, maxBits)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &costs[i])
	}
	const inf int64 = 1 << 60
	for i := n; i < maxBits; i++ {
		costs[i] = inf
	}
	// Ensure buying larger bottles is never worse than buying two smaller ones
	for i := 1; i < maxBits; i++ {
		if costs[i] > 2*costs[i-1] {
			costs[i] = 2 * costs[i-1]
		}
	}

	ans := inf
	var spent int64
	for i := maxBits - 1; i >= 0; i-- {
		size := int64(1) << uint(i)
		need := L / size
		spent += need * costs[i]
		L -= need * size
		// Option to buy an extra bottle of this size if needed
		if L > 0 {
			if candidate := spent + costs[i]; candidate < ans {
				ans = candidate
			}
		} else {
			if spent < ans {
				ans = spent
			}
		}
	}
	fmt.Fprintln(writer, ans)
}

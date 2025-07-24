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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	houses := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &houses[i])
	}

	const inf = int(1e9)
	ans := inf
	for i, cost := range houses {
		if cost != 0 && cost <= k {
			dist := abs(i+1-m) * 10
			if dist < ans {
				ans = dist
			}
		}
	}

	fmt.Fprintln(writer, ans)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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

	var k1, k2, k3 int
	if _, err := fmt.Fscan(reader, &k1, &k2, &k3); err != nil {
		return
	}
	n := k1 + k2 + k3
	owner := make([]int, n+1)
	for i := 0; i < k1; i++ {
		var x int
		fmt.Fscan(reader, &x)
		owner[x] = 1
	}
	for i := 0; i < k2; i++ {
		var x int
		fmt.Fscan(reader, &x)
		owner[x] = 2
	}
	for i := 0; i < k3; i++ {
		var x int
		fmt.Fscan(reader, &x)
		owner[x] = 3
	}

	const inf = int(1e9)
	dp1 := 0
	dp2 := inf
	dp3 := inf
	for i := 1; i <= n; i++ {
		cost1 := 0
		if owner[i] != 1 {
			cost1 = 1
		}
		cost2 := 0
		if owner[i] != 2 {
			cost2 = 1
		}
		cost3 := 0
		if owner[i] != 3 {
			cost3 = 1
		}
		newDp1 := dp1 + cost1
		newDp2 := min(dp1, dp2) + cost2
		newDp3 := min(dp2, dp3) + cost3
		dp1, dp2, dp3 = newDp1, newDp2, newDp3
	}

	fmt.Fprintln(writer, dp3)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

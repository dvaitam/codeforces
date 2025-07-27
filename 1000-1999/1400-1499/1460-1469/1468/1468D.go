package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, a, b int
		fmt.Fscan(reader, &n, &m, &a, &b)
		times := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &times[i])
		}

		var limit int
		if a < b {
			limit = b - 1
		} else {
			limit = n - b
		}
		dist := abs(a-b) - 1
		sort.Sort(sort.Reverse(sort.IntSlice(times)))
		k := min(dist, m)
		ans := 0
		for i := 0; i < k; i++ {
			if times[i]+i+1 <= limit {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

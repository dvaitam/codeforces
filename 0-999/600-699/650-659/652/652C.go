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
	perm := make([]int, n)
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &perm[i])
		pos[perm[i]] = i + 1
	}

	limit := make([]int, n+2)
	for i := range limit {
		limit[i] = n
	}

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		x := pos[a]
		y := pos[b]
		if x > y {
			x, y = y, x
		}
		if y-1 < limit[x] {
			limit[x] = y - 1
		}
	}

	for i := n - 1; i >= 1; i-- {
		if limit[i+1] < limit[i] {
			limit[i] = limit[i+1]
		}
	}

	var ans int64
	for i := 1; i <= n; i++ {
		ans += int64(limit[i] - i + 1)
	}
	fmt.Fprintln(writer, ans)
}

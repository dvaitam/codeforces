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

	var n, q int
	fmt.Fscan(reader, &n, &q)

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	parent := make([]int, n+1)
	parent[1] = 0
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &parent[i])
	}

	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		var ans int64
		for x > 0 && y > 0 {
			ans += a[x] * a[y]
			x = parent[x]
			y = parent[y]
		}
		fmt.Fprintln(writer, ans)
	}
}

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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if n < 4 {
		fmt.Fprintln(writer, 0)
		return
	}
	const INF = int(1e9)
	minX, maxX := INF, -INF
	minY, maxY := INF, -INF
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	dx := maxX - minX - 1
	dy := maxY - minY - 1
	if dx <= 0 || dy <= 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	ans := int64(dx) * int64(dy)
	fmt.Fprintln(writer, ans)
}

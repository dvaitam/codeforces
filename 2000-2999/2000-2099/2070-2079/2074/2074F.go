package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l1, r1, l2, r2 int
		fmt.Fscan(in, &l1, &r1, &l2, &r2)
		if l1 > r1 {
			l1, r1 = r1, l1
		}
		if l2 > r2 {
			l2, r2 = r2, l2
		}
		len := 1
		maxCoord := r1
		if r2 > maxCoord {
			maxCoord = r2
		}
		for len < maxCoord {
			len <<= 1
		}
		ans := countNodes(0, 0, len, l1, r1, l2, r2)
		fmt.Fprintln(out, ans)
	}
}

func countNodes(x, y, size, l1, r1, l2, r2 int) int64 {
	x2 := x + size
	y2 := y + size
	if x2 <= l1 || r1 <= x || y2 <= l2 || r2 <= y {
		return 0
	}
	if l1 <= x && x2 <= r1 && l2 <= y && y2 <= r2 {
		return 1
	}
	if size == 1 {
		return 1
	}
	half := size >> 1
	total := countNodes(x, y, half, l1, r1, l2, r2)
	total += countNodes(x+half, y, half, l1, r1, l2, r2)
	total += countNodes(x, y+half, half, l1, r1, l2, r2)
	total += countNodes(x+half, y+half, half, l1, r1, l2, r2)
	return total
}

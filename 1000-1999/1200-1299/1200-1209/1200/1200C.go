package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int64
	var q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}
	g := gcd(n, m)
	innerSeg := n / g
	outerSeg := m / g

	for ; q > 0; q-- {
		var sx, ex int
		var sy, ey int64
		fmt.Fscan(reader, &sx, &sy, &ex, &ey)
		var seg1, seg2 int64
		if sx == 1 {
			seg1 = (sy - 1) / innerSeg
		} else {
			seg1 = (sy - 1) / outerSeg
		}
		if ex == 1 {
			seg2 = (ey - 1) / innerSeg
		} else {
			seg2 = (ey - 1) / outerSeg
		}
		if seg1 == seg2 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

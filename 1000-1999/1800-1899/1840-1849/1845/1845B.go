package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var xA, yA, xB, yB, xC, yC int64
		fmt.Fscan(in, &xA, &yA)
		fmt.Fscan(in, &xB, &yB)
		fmt.Fscan(in, &xC, &yC)

		ans := int64(1)
		dx1 := xB - xA
		dx2 := xC - xA
		if dx1*dx2 > 0 {
			ans += minInt64(absInt64(dx1), absInt64(dx2))
		}
		dy1 := yB - yA
		dy2 := yC - yA
		if dy1*dy2 > 0 {
			ans += minInt64(absInt64(dy1), absInt64(dy2))
		}
		fmt.Fprintln(out, ans)
	}
}

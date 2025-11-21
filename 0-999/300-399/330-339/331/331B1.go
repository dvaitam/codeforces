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

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		pos[a[i]] = i
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(in, &t, &x, &y)
		if t == 1 {
			l := x
			r := y
			count := 1
			for v := l; v < r; v++ {
				if pos[v] > pos[v+1] {
					count++
				}
			}
			fmt.Fprintln(out, count)
		} else {
			x--
			y--
			ax, ay := a[x], a[y]
			pos[ax], pos[ay] = y, x
			a[x], a[y] = a[y], a[x]
		}
	}
}

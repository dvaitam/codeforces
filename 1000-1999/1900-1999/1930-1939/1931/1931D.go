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
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		groups := make(map[int64]map[int64]int)
		var ans int64
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(in, &a)
			ry := a % y
			rx := a % x
			want := (x - rx) % x
			if m, ok := groups[ry]; ok {
				if c, ok2 := m[want]; ok2 {
					ans += int64(c)
				}
			}
			if _, ok := groups[ry]; !ok {
				groups[ry] = make(map[int64]int)
			}
			groups[ry][rx]++
		}
		fmt.Fprintln(out, ans)
	}
}

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
		var a [4]int
		for i := 0; i < 4; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := 0
		// Check three possible pairings of segments into width and height
		area := func(i1, i2, i3, i4 int) int {
			w := a[i1]
			if a[i2] < w {
				w = a[i2]
			}
			h := a[i3]
			if a[i4] < h {
				h = a[i4]
			}
			return w * h
		}
		if v := area(0, 1, 2, 3); v > ans {
			ans = v
		}
		if v := area(0, 2, 1, 3); v > ans {
			ans = v
		}
		if v := area(0, 3, 1, 2); v > ans {
			ans = v
		}
		fmt.Fprintln(out, ans)
	}
}

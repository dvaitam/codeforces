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
		var n, x int
		fmt.Fscan(in, &n, &x)
		minPos := int(1e9)
		maxPos := -1
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v == 1 {
				if i < minPos {
					minPos = i
				}
				if i > maxPos {
					maxPos = i
				}
			}
		}
		duration := maxPos - minPos + 1
		if duration <= x {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

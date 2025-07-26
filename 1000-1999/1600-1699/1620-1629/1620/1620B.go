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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var w, h int64
		fmt.Fscan(reader, &w, &h)

		// bottom side (y=0)
		var k int
		fmt.Fscan(reader, &k)
		var first, last int64
		for i := 0; i < k; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			if i == 0 {
				first = x
			}
			if i == k-1 {
				last = x
			}
		}
		ans := (last - first) * h

		// top side (y=h)
		fmt.Fscan(reader, &k)
		for i := 0; i < k; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			if i == 0 {
				first = x
			}
			if i == k-1 {
				last = x
			}
		}
		if val := (last - first) * h; val > ans {
			ans = val
		}

		// left side (x=0)
		fmt.Fscan(reader, &k)
		for i := 0; i < k; i++ {
			var y int64
			fmt.Fscan(reader, &y)
			if i == 0 {
				first = y
			}
			if i == k-1 {
				last = y
			}
		}
		if val := (last - first) * w; val > ans {
			ans = val
		}

		// right side (x=w)
		fmt.Fscan(reader, &k)
		for i := 0; i < k; i++ {
			var y int64
			fmt.Fscan(reader, &y)
			if i == 0 {
				first = y
			}
			if i == k-1 {
				last = y
			}
		}
		if val := (last - first) * w; val > ans {
			ans = val
		}

		fmt.Fprintln(writer, ans)
	}
}

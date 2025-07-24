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
		var n int
		fmt.Fscan(reader, &n)
		maxRight, maxLeft := 0, 0
		maxUp, maxDown := 0, 0
		for i := 0; i < n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x > 0 {
				if x > maxRight {
					maxRight = x
				}
			} else if x < 0 {
				if -x > maxLeft {
					maxLeft = -x
				}
			}
			if y > 0 {
				if y > maxUp {
					maxUp = y
				}
			} else if y < 0 {
				if -y > maxDown {
					maxDown = -y
				}
			}
		}
		ans := 2 * (maxRight + maxLeft + maxUp + maxDown)
		fmt.Fprintln(writer, ans)
	}
}

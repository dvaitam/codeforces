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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var B, x, y int64
		fmt.Fscan(in, &n, &B, &x, &y)
		var cur, sum int64
		for i := 0; i < n; i++ {
			if cur+x <= B {
				cur += x
			} else {
				cur -= y
			}
			sum += cur
		}
		fmt.Fprintln(out, sum)
	}
}

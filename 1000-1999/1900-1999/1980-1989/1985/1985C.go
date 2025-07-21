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
		fmt.Fscan(in, &n)
		var sum, maxVal int64
		count := 0
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			sum += v
			if v > maxVal {
				maxVal = v
			}
			if sum == maxVal*2 {
				count++
			}
		}
		fmt.Fprintln(out, count)
	}
}

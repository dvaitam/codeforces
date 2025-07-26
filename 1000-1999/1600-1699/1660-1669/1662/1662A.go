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
		fmt.Fscan(in, &n)
		best := make([]int, 11)
		for i := 1; i <= 10; i++ {
			best[i] = -1
		}
		for i := 0; i < n; i++ {
			var b, d int
			fmt.Fscan(in, &b, &d)
			if b > best[d] {
				best[d] = b
			}
		}
		possible := true
		sum := 0
		for d := 1; d <= 10; d++ {
			if best[d] == -1 {
				possible = false
				break
			}
			sum += best[d]
		}
		if possible {
			fmt.Fprintln(out, sum)
		} else {
			fmt.Fprintln(out, "MOREPROBLEMS")
		}
	}
}

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
		if n == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		var total, best int64
		for i := 0; i < n; i++ {
			next := arr[(i+1)%n]
			edge := arr[i]
			if next > edge {
				edge = next
			}
			total += edge
			if edge > best {
				best = edge
			}
		}

		fmt.Fprintln(out, total-best)
	}
}

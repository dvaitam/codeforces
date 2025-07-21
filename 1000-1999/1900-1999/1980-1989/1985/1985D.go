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
		var n, m int
		fmt.Fscan(in, &n, &m)
		var sumX, sumY, cnt int
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var ch byte
				fmt.Fscan(in, &ch)
				if ch == '#' {
					sumX += i + 1
					sumY += j + 1
					cnt++
				}
			}
		}
		if cnt > 0 {
			fmt.Fprintf(out, "%d %d\n", sumX/cnt, sumY/cnt)
		} else {
			fmt.Fprintln(out)
		}
	}
}

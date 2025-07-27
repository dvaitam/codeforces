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
		var n, m, x, y int
		fmt.Fscan(in, &n, &m, &x, &y)
		total := 0
		for i := 0; i < n; i++ {
			var row string
			fmt.Fscan(in, &row)
			j := 0
			for j < m {
				if row[j] == '*' {
					j++
					continue
				}
				if j+1 < m && row[j+1] == '.' && y < 2*x {
					total += y
					j += 2
				} else {
					total += x
					j++
				}
			}
		}
		fmt.Fprintln(out, total)
	}
}

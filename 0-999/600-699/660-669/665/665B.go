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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	row := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &row[i])
	}

	total := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var x int
			fmt.Fscan(in, &x)
			pos := 0
			for ; pos < k; pos++ {
				if row[pos] == x {
					break
				}
			}
			total += pos + 1
			tmp := row[pos]
			for t := pos; t > 0; t-- {
				row[t] = row[t-1]
			}
			row[0] = tmp
		}
	}

	fmt.Fprintln(out, total)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		res := make([]int, n*m)
		idx := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				d1 := i + j
				d2 := i + (m - 1 - j)
				d3 := (n - 1 - i) + j
				d4 := (n - 1 - i) + (m - 1 - j)
				maxd := d1
				if d2 > maxd {
					maxd = d2
				}
				if d3 > maxd {
					maxd = d3
				}
				if d4 > maxd {
					maxd = d4
				}
				res[idx] = maxd
				idx++
			}
		}
		sort.Ints(res)
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

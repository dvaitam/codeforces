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
		var n, k int
		fmt.Fscan(in, &n, &k)
		pairs := make([][2]int, 0, n)
		used := make([]bool, 2*n+1)
		for i := 0; i < k; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x > y {
				x, y = y, x
			}
			pairs = append(pairs, [2]int{x, y})
			used[x] = true
			used[y] = true
		}
		// collect unused points
		unused := make([]int, 0, 2*(n-k))
		for i := 1; i <= 2*n; i++ {
			if !used[i] {
				unused = append(unused, i)
			}
		}
		m := n - k
		// pair first half with second half to maximize intersections
		for i := 0; i < m; i++ {
			x := unused[i]
			y := unused[i+m]
			if x > y {
				x, y = y, x
			}
			pairs = append(pairs, [2]int{x, y})
		}
		// count intersections
		sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
		count := 0
		total := len(pairs)
		for i := 0; i < total; i++ {
			a, b := pairs[i][0], pairs[i][1]
			for j := i + 1; j < total; j++ {
				c, d := pairs[j][0], pairs[j][1]
				if a < c && c < b && b < d {
					count++
				} else if c < a && a < d && d < b {
					count++
				}
			}
		}
		fmt.Fprintln(out, count)
	}
}

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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, l, r int
		fmt.Fscan(in, &n, &l, &r)
		inside := make([]int, 0, r-l+1)
		outside := make([]int, 0, n-(r-l+1))
		var sum int64
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if l <= i && i <= r {
				inside = append(inside, x)
				sum += int64(x)
			} else {
				outside = append(outside, x)
			}
		}
		sort.Slice(inside, func(i, j int) bool { return inside[i] > inside[j] })
		sort.Ints(outside)
		limit := len(inside)
		if len(outside) < limit {
			limit = len(outside)
		}
		for i := 0; i < limit; i++ {
			if outside[i] < inside[i] {
				sum -= int64(inside[i] - outside[i])
			} else {
				break
			}
		}
		fmt.Fprintln(out, sum)
	}
}

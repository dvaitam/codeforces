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
		br := make([][]int, k+1)
		for i := 0; i < k; i++ {
			var b, c int
			fmt.Fscan(in, &b, &c)
			br[b] = append(br[b], c)
		}
		sums := make([]int, 0)
		for _, costs := range br {
			if len(costs) == 0 {
				continue
			}
			sort.Sort(sort.Reverse(sort.IntSlice(costs)))
			prefix := 0
			for i := 0; i < len(costs); i++ {
				prefix += costs[i]
				sums = append(sums, prefix)
				if i+1 >= n {
					break
				}
			}
		}
		sort.Sort(sort.Reverse(sort.IntSlice(sums)))
		res := 0
		for i := 0; i < len(sums) && i < n; i++ {
			res += sums[i]
		}
		fmt.Fprintln(out, res)
	}
}

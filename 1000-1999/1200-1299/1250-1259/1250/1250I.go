package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Subset struct {
	size int
	sum  int64
	mask uint64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		var k int64
		fmt.Fscan(in, &n, &k, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		subsets := make([]Subset, 0)
		var dfs func(idx int, sum int64, mask uint64, size int)
		dfs = func(idx int, sum int64, mask uint64, size int) {
			if sum > k {
				return
			}
			if idx == n {
				if size > 0 {
					subsets = append(subsets, Subset{size, sum, mask})
				}
				return
			}
			dfs(idx+1, sum, mask, size)
			dfs(idx+1, sum+a[idx], mask|1<<uint(idx), size+1)
		}
		dfs(0, 0, 0, 0)

		sort.Slice(subsets, func(i, j int) bool {
			if subsets[i].size != subsets[j].size {
				return subsets[i].size > subsets[j].size
			}
			if subsets[i].sum != subsets[j].sum {
				return subsets[i].sum < subsets[j].sum
			}
			return subsets[i].mask < subsets[j].mask
		})

		r := m
		if r > len(subsets) {
			r = len(subsets)
		}
		if r == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		fmt.Fprintln(out, r)
		for i := 0; i < r; i++ {
			fmt.Fprintf(out, "%d %d\n", subsets[i].size, subsets[i].sum)
		}
		last := subsets[r-1]
		first := true
		for i := 0; i < n; i++ {
			if last.mask&(1<<uint(i)) != 0 {
				if !first {
					out.WriteByte(' ')
				}
				first = false
				fmt.Fprint(out, i+1)
			}
		}
		if !first {
			out.WriteByte('\n')
		} else {
			fmt.Fprintln(out)
		}
	}
}

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
		var n int64
		var m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, m)
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}
		idx := make([]int, m)
		for i := 0; i < m; i++ {
			idx[i] = i
		}
		sort.Slice(idx, func(i, j int) bool { return a[idx[i]] > a[idx[j]] })
		sortedA := make([]int64, m)
		pos := make([]int, m)
		for i := 0; i < m; i++ {
			sortedA[i] = a[idx[i]]
			pos[idx[i]] = i
		}
		prefix := make([]int64, m+1)
		for i := 0; i < m; i++ {
			prefix[i+1] = prefix[i] + sortedA[i]
		}
		best := int64(0)
		for i := 0; i < m; i++ {
			bi := b[i]
			k := sort.Search(m, func(p int) bool { return sortedA[p] <= bi })
			if int64(k) > n {
				k = int(n)
			}
			used := int64(k)
			res := prefix[k]
			if pos[i] >= k {
				if used < n {
					res += a[i]
					used++
				} else {
					if k > 0 {
						res = res - sortedA[k-1] + a[i]
					} else {
						res += a[i]
					}
				}
			}
			if used < n {
				res += (n - used) * bi
			}
			if res > best {
				best = res
			}
		}
		fmt.Fprintln(out, best)
	}
}

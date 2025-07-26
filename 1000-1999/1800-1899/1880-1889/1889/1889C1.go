package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		starts := make([][]int, n+2)
		ends := make([][]int, n+2)
		l := make([]int, m+1)
		r := make([]int, m+1)
		for i := 1; i <= m; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
			starts[l[i]] = append(starts[l[i]], i)
			ends[r[i]] = append(ends[r[i]], i)
		}
		active := make(map[int]struct{})
		single := make([]int, m+1)
		pairCnt := make(map[[2]int]int)
		base := 0
		for pos := 1; pos <= n; pos++ {
			for _, id := range starts[pos] {
				active[id] = struct{}{}
			}
			switch len(active) {
			case 0:
				base++
			case 1:
				for id := range active {
					single[id]++
				}
			case 2:
				ids := make([]int, 0, 2)
				for id := range active {
					ids = append(ids, id)
				}
				a, b := ids[0], ids[1]
				if a > b {
					a, b = b, a
				}
				pairCnt[[2]int{a, b}]++
			}
			for _, id := range ends[pos] {
				delete(active, id)
			}
		}
		max1, max2 := 0, 0
		for i := 1; i <= m; i++ {
			if single[i] > max1 {
				max2 = max1
				max1 = single[i]
			} else if single[i] > max2 {
				max2 = single[i]
			}
		}
		ans := base + max1 + max2
		for key, cnt := range pairCnt {
			cand := base + single[key[0]] + single[key[1]] + cnt
			if cand > ans {
				ans = cand
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Album struct {
	max  int
	pref []int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		albums := make([]Album, n)
		vals := []int{0}
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(in, &k)
			arr := make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(in, &arr[j])
			}
			pref := make([]int, 0, k)
			cur := 0
			for _, x := range arr {
				if x > cur {
					cur = x
					pref = append(pref, x)
				}
			}
			albums[i] = Album{max: cur, pref: pref}
			vals = append(vals, cur)
		}
		sort.Ints(vals)
		uniq := make([]int, 0, len(vals))
		last := -1
		for _, v := range vals {
			if v != last {
				uniq = append(uniq, v)
				last = v
			}
		}
		idx := make(map[int]int, len(uniq))
		for i, v := range uniq {
			idx[v] = i
		}
		groups := make([][]Album, len(uniq))
		for _, al := range albums {
			groups[idx[al.max]] = append(groups[idx[al.max]], al)
		}
		dp := make([]int, len(uniq))
		for i := range dp {
			dp[i] = -1 << 60
		}
		dp[0] = 0
		prefix := make([]int, len(uniq))
		prefix[0] = 0
		for id := 1; id < len(uniq); id++ {
			best := dp[id]
			for _, al := range groups[id] {
				m := len(al.pref)
				for j, p := range al.pref {
					pos := sort.SearchInts(uniq, p) - 1
					cand := m - j
					if pos >= 0 {
						if prefix[pos]+cand > best {
							best = prefix[pos] + cand
						}
					} else if cand > best {
						best = cand
					}
				}
			}
			if best > dp[id] {
				dp[id] = best
			}
			if dp[id] > prefix[id-1] {
				prefix[id] = dp[id]
			} else {
				prefix[id] = prefix[id-1]
			}
		}
		fmt.Fprintln(out, prefix[len(prefix)-1])
	}
}

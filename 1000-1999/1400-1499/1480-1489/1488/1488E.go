package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ l, r int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		pos := make(map[int][]int)
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(in, &x)
			pos[x] = append(pos[x], i)
		}
		pairs := make([]pair, 0)
		for _, p := range pos {
			if len(p) == 2 {
				if p[0] > p[1] {
					p[0], p[1] = p[1], p[0]
				}
				pairs = append(pairs, pair{p[0], p[1]})
			}
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].l == pairs[j].l {
				return pairs[i].r > pairs[j].r
			}
			return pairs[i].l < pairs[j].l
		})
		dpVals := make([]int, 0)
		lens := make([]int, len(pairs))
		for idx, p := range pairs {
			v := -p.r
			pos := sort.Search(len(dpVals), func(i int) bool { return dpVals[i] >= v })
			if pos == len(dpVals) {
				dpVals = append(dpVals, v)
			} else {
				dpVals[pos] = v
			}
			lens[idx] = pos + 1
		}
		pairCnt := len(dpVals)
		extra := false
		for i, p := range pairs {
			if lens[i] == pairCnt && p.r-p.l > 1 {
				extra = true
				break
			}
		}
		res := 0
		if pairCnt == 0 {
			if n > 0 {
				res = 1
			}
		} else {
			res = pairCnt * 2
			if extra {
				res++
			}
		}
		fmt.Fprintln(out, res)
	}
}

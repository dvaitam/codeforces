package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segment struct {
	l   int64
	r   int64
	sum int64
	idx int
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

		segs := make([]segment, n)
		var origSum int64
		var totalL int64

		for i := 0; i < n; i++ {
			var l, r int64
			fmt.Fscan(in, &l, &r)
			segs[i] = segment{l: l, r: r, sum: l + r, idx: i}
			origSum += r - l
			totalL += l
		}

		order := make([]segment, n)
		copy(order, segs)
		sort.Slice(order, func(i, j int) bool {
			if order[i].sum == order[j].sum {
				if order[i].l == order[j].l {
					return order[i].idx < order[j].idx
				}
				return order[i].l < order[j].l
			}
			return order[i].sum > order[j].sum
		})

		k := n / 2
		var sumTop int64
		inTop := make([]bool, n)
		for i := 0; i < k; i++ {
			sumTop += order[i].sum
			inTop[order[i].idx] = true
		}

		var extra int64
		if n%2 == 0 {
			extra = -totalL + sumTop
		} else {
			nextVal := order[k].sum
			var maxCand int64
			for i := 0; i < n; i++ {
				top := sumTop
				if inTop[i] {
					top = sumTop - segs[i].sum + nextVal
				}
				cand := top + segs[i].l
				if cand > maxCand {
					maxCand = cand
				}
			}
			extra = -totalL + maxCand
		}

		ans := origSum + extra
		fmt.Fprintln(out, ans)
	}
}

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

	var n, k int
	fmt.Fscan(in, &n, &k)
	type seg struct{ r, l, idx int }
	s := make([]seg, n)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		s[i] = seg{r: r, l: l, idx: i + 1}
	}
	sort.Slice(s, func(i, j int) bool {
		if s[i].r != s[j].r {
			return s[i].r < s[j].r
		}
		if s[i].l != s[j].l {
			return s[i].l < s[j].l
		}
		return s[i].idx < s[j].idx
	})

	x := make([]int, 201)
	var removed []int
	for _, sg := range s {
		ok := true
		for j := sg.l; j <= sg.r; j++ {
			if x[j]+1 > k {
				ok = false
				break
			}
		}
		if ok {
			for j := sg.l; j <= sg.r; j++ {
				x[j]++
			}
		} else {
			removed = append(removed, sg.idx)
		}
	}

	fmt.Fprintln(out, len(removed))
	for _, idx := range removed {
		fmt.Fprint(out, idx, " ")
	}
	fmt.Fprintln(out)
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type segment struct {
	l int
	r int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return
		}
		segs := make([]segment, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
		}
		var q int
		fmt.Fscan(in, &q)
		qpos := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &qpos[i])
		}

		arr := make([]int, n+1)
		pref := make([]int, n+1)
		check := func(k int) bool {
			for i := 0; i <= n; i++ {
				arr[i] = 0
			}
			for i := 0; i < k; i++ {
				arr[qpos[i]] = 1
			}
			pref[0] = 0
			for i := 1; i <= n; i++ {
				pref[i] = pref[i-1] + arr[i]
			}
			for _, s := range segs {
				ones := pref[s.r] - pref[s.l-1]
				if ones*2 > s.r-s.l+1 {
					return true
				}
			}
			return false
		}

		if !check(q) {
			fmt.Fprintln(out, -1)
			continue
		}
		l, r := 1, q
		ans := q
		for l <= r {
			mid := (l + r) / 2
			if check(mid) {
				ans = mid
				r = mid - 1
			} else {
				l = mid + 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}

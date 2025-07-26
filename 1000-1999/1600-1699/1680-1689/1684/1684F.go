package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct{ l, r int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		mx := make([]int, n+2)
		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			if r > mx[l] {
				mx[l] = r
			}
		}
		for i := 1; i <= n; i++ {
			if mx[i] < mx[i-1] {
				mx[i] = mx[i-1]
			}
		}
		pos := make(map[int][]int)
		for i := 1; i <= n; i++ {
			pos[a[i]] = append(pos[a[i]], i)
		}
		pairs := make([]Pair, 0)
		for _, v := range pos {
			for j := 0; j+1 < len(v); j++ {
				x := v[j]
				y := v[j+1]
				if mx[x] >= y {
					pairs = append(pairs, Pair{x, y})
				}
			}
		}
		if len(pairs) == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		endpoints := make([][]int, n+2)
		for idx, p := range pairs {
			endpoints[p.l] = append(endpoints[p.l], idx)
			endpoints[p.r] = append(endpoints[p.r], idx)
		}
		count := make([]int, len(pairs))
		covered := 0
		ans := n
		R := 0
		for L := 1; L <= n; L++ {
			for R < n && covered < len(pairs) {
				R++
				for _, id := range endpoints[R] {
					if count[id] == 0 {
						covered++
					}
					count[id]++
				}
			}
			if covered == len(pairs) {
				if R-L+1 < ans {
					ans = R - L + 1
				}
			} else {
				break
			}
			for _, id := range endpoints[L] {
				count[id]--
				if count[id] == 0 {
					covered--
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

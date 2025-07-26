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
		var n, m int
		fmt.Fscan(in, &n, &m)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
		}
		lcnt := 0
		rcnt := 0
		posMap := make(map[int]bool)
		for _, v := range vals {
			if v == -1 {
				lcnt++
			} else if v == -2 {
				rcnt++
			} else if v > 0 {
				posMap[v] = true
			}
		}
		pos := make([]int, 0, len(posMap))
		for k := range posMap {
			if k >= 1 && k <= m {
				pos = append(pos, k)
			}
		}
		sort.Ints(pos)
		k := len(pos)
		ans := 0
		// case: fill left or right without considering specific pivot
		if lcnt+k > ans {
			tmp := lcnt + k
			if tmp > m {
				tmp = m
			}
			if tmp > ans {
				ans = tmp
			}
		}
		if rcnt+k > ans {
			tmp := rcnt + k
			if tmp > m {
				tmp = m
			}
			if tmp > ans {
				ans = tmp
			}
		}
		// consider each pivot seat
		for i, p := range pos {
			left := p - 1
			if left > lcnt+i {
				left = lcnt + i
			}
			right := m - p
			if right > rcnt+(k-i-1) {
				right = rcnt + (k - i - 1)
			}
			cur := 1 + left + right
			if cur > ans {
				ans = cur
			}
		}
		if ans > m {
			ans = m
		}
		fmt.Fprintln(out, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countGE(arr []int, x int) int64 {
	idx := sort.SearchInts(arr, x)
	return int64(len(arr) - idx)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for tc := 0; tc < T; tc++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, m)
		for i := range b {
			fmt.Fscan(in, &b[i])
		}
		kevin := a[0]
		aSorted := append([]int(nil), a...)
		sort.Ints(aSorted)

		hard := make([]int, 0)
		easyCount := 0
		for _, val := range b {
			if val <= kevin {
				easyCount++
			} else {
				hard = append(hard, val)
			}
		}
		sort.Ints(hard)
		hardLen := len(hard)

		for k := 1; k <= m; k++ {
			contests := m / k
			if contests == 0 {
				fmt.Fprint(out, 0)
				if k < m {
					fmt.Fprint(out, " ")
				}
				continue
			}
			selected := contests * k
			cnt := selected - easyCount
			if cnt < 0 {
				cnt = 0
			}
			if cnt > hardLen {
				cnt = hardLen
			}
			base := int64(contests)
			if cnt == 0 {
				fmt.Fprint(out, base)
				if k < m {
					fmt.Fprint(out, " ")
				}
				continue
			}
			start := hardLen - cnt
			groups := (cnt + k - 1) / k
			penalty := int64(0)
			for g := 0; g < groups; g++ {
				idx := start + g*k
				thr := hard[idx]
				penalty += countGE(aSorted, thr)
			}
			ans := base + penalty
			fmt.Fprint(out, ans)
			if k < m {
				fmt.Fprint(out, " ")
			}
		}
		if tc+1 < T {
			fmt.Fprintln(out)
		}
	}
}

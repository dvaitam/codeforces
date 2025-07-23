package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	diff := make([]int, m+2)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		if l < 1 {
			l = 1
		}
		if r > m {
			r = m
		}
		diff[l]++
		if r+1 <= m {
			diff[r+1]--
		}
	}
	cnt := make([]int, m+1)
	cur := 0
	for i := 1; i <= m; i++ {
		cur += diff[i]
		cnt[i] = cur
	}

	lis := make([]int, m+1)
	tails := []int{}
	for i := 1; i <= m; i++ {
		x := cnt[i]
		idx := sort.Search(len(tails), func(j int) bool { return tails[j] > x })
		if idx == len(tails) {
			tails = append(tails, x)
		} else {
			tails[idx] = x
		}
		lis[i] = idx + 1
	}

	lds := make([]int, m+2)
	tails = tails[:0]
	for i := m; i >= 1; i-- {
		x := cnt[i]
		idx := sort.Search(len(tails), func(j int) bool { return tails[j] > x })
		if idx == len(tails) {
			tails = append(tails, x)
		} else {
			tails[idx] = x
		}
		lds[i] = idx + 1
	}

	ans := 0
	for i := 1; i <= m; i++ {
		if v := lis[i] + lds[i] - 1; v > ans {
			ans = v
		}
	}
	fmt.Println(ans)
}

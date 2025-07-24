package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cup struct {
	c int64
	w int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	var d int64
	if _, err := fmt.Fscan(in, &n, &m, &d); err != nil {
		return
	}
	phys := make([]Cup, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &phys[i].c, &phys[i].w)
	}
	info := make([]Cup, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &info[i].c, &info[i].w)
	}

	sort.Slice(phys, func(i, j int) bool { return phys[i].c > phys[j].c })
	sort.Slice(info, func(i, j int) bool { return info[i].c > info[j].c })

	wp := make([]int64, n+1)
	sp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		wp[i] = wp[i-1] + phys[i-1].w
		sp[i] = sp[i-1] + phys[i-1].c
	}
	wi := make([]int64, m+1)
	si := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		wi[i] = wi[i-1] + info[i-1].w
		si[i] = si[i-1] + info[i-1].c
	}

	ans := int64(0)
	for i := 1; i <= n; i++ {
		remain := d - wp[i]
		if remain < 0 {
			continue
		}
		j := sort.Search(len(wi), func(k int) bool { return wi[k] > remain }) - 1
		if j >= 1 {
			val := sp[i] + si[j]
			if val > ans {
				ans = val
			}
		}
	}
	for j := 1; j <= m; j++ {
		remain := d - wi[j]
		if remain < 0 {
			continue
		}
		i := sort.Search(len(wp), func(k int) bool { return wp[k] > remain }) - 1
		if i >= 1 {
			val := si[j] + sp[i]
			if val > ans {
				ans = val
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprint(out, ans)
	out.Flush()
}

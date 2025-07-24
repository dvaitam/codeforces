package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	a   int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{a: a[i], idx: i + 1}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].a < arr[j].a })
		prefix := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + arr[i].a
		}
		k := 0
		for k < n && prefix[k+1] <= m {
			k++
		}
		if k == n {
			fmt.Fprintln(out, 1)
			continue
		}
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			pos[arr[i].idx] = i + 1 // 1-indexed position
		}
		costInclude := prefix[k]
		p := pos[k+1]
		if p > k {
			if k == 0 {
				costInclude = arr[p-1].a
			} else {
				costInclude = prefix[k] - arr[k-1].a + arr[p-1].a
			}
		}
		if costInclude <= m {
			fmt.Fprintln(out, n-k)
		} else {
			fmt.Fprintln(out, n-k+1)
		}
	}
}

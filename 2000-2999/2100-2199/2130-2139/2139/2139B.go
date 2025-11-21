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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		limit := n
		if int64(limit) > m {
			limit = int(m)
		}
		var ans int64
		for i := 0; i < limit; i++ {
			ans += a[i] * (m - int64(i))
		}
		fmt.Fprintln(out, ans)
	}
}

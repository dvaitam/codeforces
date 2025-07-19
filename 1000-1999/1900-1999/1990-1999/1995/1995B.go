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
		v := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &v[i])
		}
		sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
		l, r := 0, 1
		sum := v[0]
		var ans int64
		if v[0] > m {
			fmt.Fprintln(out, 0)
			continue
		}
		for l < n {
			if r < n && v[r]-v[l] <= 1 && sum+v[r] <= m {
				sum += v[r]
				r++
			} else {
				sum -= v[l]
				l++
			}
			if sum > ans {
				ans = sum
			}
		}
		if v[0] > ans {
			ans = v[0]
		}
		fmt.Fprintln(out, ans)
	}
}

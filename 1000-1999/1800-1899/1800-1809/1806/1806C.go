package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int64, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(in, &p[i])
		}
		sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
		if n == 1 {
			ans := absInt64(p[0] - p[1])
			fmt.Fprintln(out, ans)
			continue
		}
		const inf int64 = 1 << 60
		ans := inf
		if n == 2 {
			var tmp int64
			for i := 0; i < 4; i++ {
				tmp += absInt64(p[i] - 2)
			}
			if tmp < ans {
				ans = tmp
			}
		}
		var sum int64
		for _, v := range p {
			sum += absInt64(v)
		}
		if sum < ans {
			ans = sum
		}
		if n%2 == 0 {
			var tmp int64
			for i := 0; i < 2*n-1; i++ {
				tmp += absInt64(p[i] + 1)
			}
			tmp += absInt64(int64(n) - p[2*n-1])
			if tmp < ans {
				ans = tmp
			}
		}
		fmt.Fprintln(out, ans)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(x int64, n int, a, b []int64) (int64, int64) {
	var cnt, ans int64
	for i := 1; i <= n; i++ {
		if a[i] < x {
			continue
		}
		now := (a[i]-x)/b[i] + 1
		cnt += now
		ans += (a[i] + a[i] - b[i]*(now-1)) * now / 2
	}
	return cnt, ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		l, r := int64(0), int64(1e9)
		for l < r {
			mid := (l + r + 1) / 2
			cnt, _ := calc(mid, n, a, b)
			if cnt >= k {
				l = mid
			} else {
				r = mid - 1
			}
		}
		cnt, ans := calc(l, n, a, b)
		fmt.Fprintln(writer, ans-(cnt-k)*l)
	}
}

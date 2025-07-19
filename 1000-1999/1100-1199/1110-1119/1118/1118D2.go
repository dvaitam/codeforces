package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	// check if using X rounds yields at least m energy
	check := func(X int) bool {
		var tmp int64
		cnt := 0
		cur := int64(0)
		for i := n - 1; i >= 0; i-- {
			if a[i] > cur {
				tmp += a[i] - cur
			}
			cnt++
			if cnt == X {
				cur++
				cnt = 0
			}
		}
		return tmp >= m
	}

	l, r := 1, n
	ans := -1
	for l <= r {
		mid := (l + r) / 2
		if check(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	fmt.Fprintln(writer, ans)
}

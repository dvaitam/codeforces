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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		l := make([]int, n)
		r := make([]int, n)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &l[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &r[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &c[i])
		}

		sort.Ints(l)
		sort.Ints(r)
		sort.Slice(c, func(i, j int) bool { return c[i] > c[j] })

		diff := make([]int, n)
		pool := make([]int, 0, n)
		idx := 0
		for i, rv := range r {
			for idx < n && l[idx] < rv {
				pool = append(pool, l[idx])
				idx++
			}
			// select largest available left < rv
			lv := pool[len(pool)-1]
			pool = pool[:len(pool)-1]
			diff[i] = rv - lv
		}
		sort.Ints(diff)

		var ans int64
		for i := 0; i < n; i++ {
			ans += int64(diff[i]) * int64(c[i])
		}
		fmt.Fprintln(writer, ans)
	}
}

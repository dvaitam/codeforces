package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(time int64, cnt []int64, m int64) bool {
	var total int64
	for _, c := range cnt {
		if c >= time {
			total += time
		} else {
			total += c + (time-c)/2
		}
		if total >= m {
			return true
		}
	}
	return total >= m
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		cnt := make([]int64, n)
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(reader, &x)
			cnt[x-1]++
		}
		l, r := int64(1), int64(2*m)
		for l < r {
			mid := (l + r) / 2
			if can(mid, cnt, int64(m)) {
				r = mid
			} else {
				l = mid + 1
			}
		}
		fmt.Fprintln(writer, l)
	}
}

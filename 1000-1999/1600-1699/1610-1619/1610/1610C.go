package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(k int, a, b []int) bool {
	cnt := 0
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] >= cnt && b[i] >= k-1-cnt {
			cnt++
			if cnt == k {
				return true
			}
		}
	}
	return cnt >= k
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i], &b[i])
		}
		l, r := 0, n
		for l < r {
			mid := (l + r + 1) / 2
			if can(mid, a, b) {
				l = mid
			} else {
				r = mid - 1
			}
		}
		fmt.Fprintln(writer, l)
	}
}

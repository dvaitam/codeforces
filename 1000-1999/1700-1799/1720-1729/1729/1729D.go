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
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &x[i])
		}
		y := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &y[i])
		}
		diff := make([]int, n)
		for i := 0; i < n; i++ {
			diff[i] = y[i] - x[i]
		}
		sort.Ints(diff)
		l, r := 0, n-1
		ans := 0
		for l < r && diff[l] < 0 {
			if diff[l]+diff[r] >= 0 {
				ans++
				l++
				r--
			} else {
				l++
			}
		}
		if l < r {
			ans += (r - l + 1) / 2
		}
		fmt.Fprintln(writer, ans)
	}
}

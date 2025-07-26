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
		var l, r, k int64
		fmt.Fscan(reader, &n, &l, &r, &k)
		prices := make([]int, 0, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			if x >= l && x <= r {
				prices = append(prices, int(x))
			}
		}
		sort.Ints(prices)
		count := 0
		budget := k
		for _, p := range prices {
			if int64(p) <= budget {
				budget -= int64(p)
				count++
			} else {
				break
			}
		}
		fmt.Fprintln(writer, count)
	}
}

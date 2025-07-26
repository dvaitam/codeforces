package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// This program solves the problem described in problemB.txt for contest 1697.
// We sort item prices in descending order and build prefix sums. For each
// query (x, y) the customer chooses the x most expensive items and receives the
// y cheapest among them for free, which are exactly the last y items among those
// x. Their total value is prefix[x] - prefix[x-y].
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	prices := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &prices[i])
	}
	sort.Slice(prices, func(i, j int) bool { return prices[i] > prices[j] })

	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(prices[i])
	}

	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		ans := prefix[x] - prefix[x-y]
		fmt.Fprintln(writer, ans)
	}
}

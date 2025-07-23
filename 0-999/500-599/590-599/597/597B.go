package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// interval represents a rental request with start and finish times.
type interval struct {
	l int
	r int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	orders := make([]interval, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &orders[i].l, &orders[i].r)
	}

	// Sort orders by finish time to apply a greedy algorithm.
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].r == orders[j].r {
			return orders[i].l < orders[j].l
		}
		return orders[i].r < orders[j].r
	})

	count := 0
	lastEnd := 0
	for _, ord := range orders {
		// Accept the order only if it starts strictly after the last accepted one.
		if ord.l > lastEnd {
			count++
			lastEnd = ord.r
		}
	}

	fmt.Fprintln(writer, count)
}

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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)

	queue := make([]int, 0, n)
	left := 0
	removed := 0
	for _, t := range a {
		limit := t - m + 1
		for left < len(queue) && queue[left] < limit {
			left++
		}
		queue = append(queue, t)
		if len(queue)-left >= k {
			queue = queue[:len(queue)-1]
			removed++
		}
	}

	fmt.Fprintln(writer, removed)
}

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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		diffs := make([]int, n-1)
		total := 0
		for i := 0; i < n-1; i++ {
			d := a[i+1] - a[i]
			if d < 0 {
				d = -d
			}
			diffs[i] = d
			total += d
		}
		sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
		remove := 0
		for i := 0; i < k-1 && i < len(diffs); i++ {
			remove += diffs[i]
		}
		fmt.Fprintln(writer, total-remove)
	}
}

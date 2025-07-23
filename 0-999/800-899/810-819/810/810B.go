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

	var n, f int
	if _, err := fmt.Fscan(reader, &n, &f); err != nil {
		return
	}
	diffs := make([]int64, n)
	var base int64
	for i := 0; i < n; i++ {
		var k, l int64
		fmt.Fscan(reader, &k, &l)
		sold := k
		if l < k {
			sold = l
		}
		base += sold
		// potential additional sales if we double k
		doubleK := k * 2
		doubleSold := doubleK
		if l < doubleK {
			doubleSold = l
		}
		diffs[i] = doubleSold - sold
	}
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i] > diffs[j]
	})
	var add int64
	if f > n {
		f = n
	}
	for i := 0; i < f; i++ {
		if diffs[i] > 0 {
			add += diffs[i]
		} else {
			break
		}
	}
	fmt.Fprintln(writer, base+add)
}

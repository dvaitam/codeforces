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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		sort.Ints(a)
		sort.Ints(b)
		best := 0
		for shift := 0; shift < n; shift++ {
			minDiff := int(^uint(0) >> 1) // int max
			for i := 0; i < n; i++ {
				diff := a[i] - b[(i+shift)%n]
				if diff < 0 {
					diff = -diff
				}
				if diff < minDiff {
					minDiff = diff
					if minDiff <= best {
						break
					}
				}
			}
			if minDiff > best {
				best = minDiff
			}
		}
		fmt.Fprintln(writer, best)
	}
}

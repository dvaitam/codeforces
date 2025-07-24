package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	minDiff := a[1] - a[0]
	count := 1
	for i := 1; i < n-1; i++ {
		diff := a[i+1] - a[i]
		if diff < minDiff {
			minDiff = diff
			count = 1
		} else if diff == minDiff {
			count++
		}
	}
	fmt.Printf("%d %d\n", minDiff, count)
}

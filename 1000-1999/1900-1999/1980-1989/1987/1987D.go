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
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)
		prev := -1
		left, right := 0, n-1
		eaten := 0
		for left <= right {
			for left <= right && a[left] <= prev {
				left++
			}
			if left > right {
				break
			}
			prev = a[left]
			left++
			eaten++
			if left <= right {
				right--
			}
		}
		fmt.Fprintln(writer, eaten)
	}
}

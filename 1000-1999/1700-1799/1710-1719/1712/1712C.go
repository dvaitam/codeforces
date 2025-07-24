package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(a []int) int {
	n := len(a)
	// last occurrence of each value (values up to n)
	last := make([]int, n+1)
	for i, v := range a {
		last[v] = i
	}
	// find first index from the right where array decreases
	end := -1
	for i := n - 2; i >= 0; i-- {
		if a[i] > a[i+1] {
			end = i
			break
		}
	}
	if end == -1 {
		return 0
	}
	seen := make(map[int]struct{})
	for i := 0; i <= end; i++ {
		v := a[i]
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			if last[v] > end {
				end = last[v]
			}
		}
	}
	return len(seen)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, solve(arr))
	}
}

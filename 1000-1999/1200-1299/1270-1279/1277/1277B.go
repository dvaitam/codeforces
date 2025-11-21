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

		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		seen := make(map[int]bool)
		moves := 0

		for _, v := range a {
			for v%2 == 0 && !seen[v] {
				seen[v] = true
				moves++
				v /= 2
			}
		}

		fmt.Fprintln(writer, moves)
	}
}

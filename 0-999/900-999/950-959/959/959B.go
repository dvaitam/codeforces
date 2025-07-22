package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, m int
	if _, err := fmt.Fscan(reader, &n, &k, &m); err != nil {
		return
	}

	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}

	costs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &costs[i])
	}

	wordIndex := make(map[string]int, n)
	for i, w := range words {
		wordIndex[w] = i
	}

	group := make([]int, n)
	minCost := make([]int64, k)
	for g := 0; g < k; g++ {
		var x int
		fmt.Fscan(reader, &x)
		minVal := int64(1<<63 - 1)
		for j := 0; j < x; j++ {
			var idx int
			fmt.Fscan(reader, &idx)
			idx--
			group[idx] = g
			if costs[idx] < minVal {
				minVal = costs[idx]
			}
		}
		minCost[g] = minVal
	}

	var total int64
	for i := 0; i < m; i++ {
		var w string
		fmt.Fscan(reader, &w)
		idx := wordIndex[w]
		g := group[idx]
		total += minCost[g]
	}

	fmt.Fprintln(writer, total)
}

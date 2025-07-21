package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Given n price tags and m fruits in the shopping list (with possible repeats),
// compute the minimum and maximum total cost by optimally assigning price tags to fruit types.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	prices := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &prices[i])
	}

	// Count each fruit in Valera's list
	countMap := make(map[string]int)
	for i := 0; i < m; i++ {
		var name string
		fmt.Fscan(reader, &name)
		countMap[name]++
	}

	// Extract counts and sort descending
	counts := make([]int, 0, len(countMap))
	for _, c := range countMap {
		counts = append(counts, c)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	// Sort prices ascending
	sort.Ints(prices)
	// Number of distinct fruits in the list
	k := len(counts)

	// Minimum sum: assign smallest prices to largest counts
	minSum := 0
	for i := 0; i < k; i++ {
		minSum += counts[i] * prices[i]
	}

	// Maximum sum: assign largest prices to largest counts
	maxSum := 0
	for i := 0; i < k; i++ {
		maxSum += counts[i] * prices[n-1-i]
	}

	fmt.Fprintln(writer, minSum, maxSum)
}

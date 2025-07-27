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
	fmt.Fscan(reader, &n)

	// Precompute proper divisors for each number.
	divs := make([][]int, n+1)
	for d := 1; d <= n; d++ {
		for m := d * 2; m <= n; m += d {
			divs[m] = append(divs[m], d)
		}
	}

	// Sort numbers by divisor count (descending), then by value.
	idx := make([]int, n)
	for i := 1; i <= n; i++ {
		idx[i-1] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		a, b := idx[i], idx[j]
		if len(divs[a]) != len(divs[b]) {
			return len(divs[a]) > len(divs[b])
		}
		return a > b
	})

	selected := make([]bool, n+1)
	countMul := make([]int, n+1)
	sumDiv := 0
	edges := 0
	ans := make([]int, n+1)

	for i, v := range idx {
		// Pairs with previously selected numbers that are multiples of v
		edges += countMul[v]
		// Pairs with previously selected divisors of v
		for _, d := range divs[v] {
			if selected[d] {
				edges++
			}
			countMul[d]++
		}
		// Mark v as selected and update multiple counts
		countMul[v]++
		selected[v] = true
		sumDiv += len(divs[v])
		ans[i+1] = sumDiv - edges
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for k := 1; k <= n; k++ {
		if k > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[k])
	}
	writer.WriteByte('\n')
}

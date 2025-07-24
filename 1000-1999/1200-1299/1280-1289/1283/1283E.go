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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	xs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i])
	}
	sort.Ints(xs)

	// Minimum number of occupied houses
	minAns := 0
	for i := 0; i < n; {
		minAns++
		limit := xs[i] + 2
		for i < n && xs[i] <= limit {
			i++
		}
	}

	// Maximum number of occupied houses
	used := make([]bool, n+3) // positions 0..n+2 possible
	for _, x := range xs {
		if x-1 >= 0 && !used[x-1] {
			used[x-1] = true
		} else if !used[x] {
			used[x] = true
		} else if !used[x+1] {
			used[x+1] = true
		}
	}
	maxAns := 0
	for _, u := range used {
		if u {
			maxAns++
		}
	}

	fmt.Fprintln(writer, minAns, maxAns)
}

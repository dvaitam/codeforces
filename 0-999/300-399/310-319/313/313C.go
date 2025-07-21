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
	// Read all elements
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// Sum of all elements
	var sumAll int64
	for _, v := range a {
		sumAll += v
	}
	// Number of internal nodes = (n - 1) / 3
	// n = 4^k, so (n-1) divisible by 3
	m := (n - 1) / 3
	// Sort ascending to take largest m elements
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var sumMax int64
	for i := n - m; i < n; i++ {
		sumMax += a[i]
	}
	// Result is sum of all elements plus sum of internal maxima
	fmt.Fprintln(writer, sumAll+sumMax)
}

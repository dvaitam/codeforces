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
	a := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
	}
	// Check first and last elements
	if a[0] != c[0] || a[n-1] != c[n-1] {
		fmt.Fprint(writer, "No")
		return
	}
	// Compute difference arrays in-place
	for i := n - 1; i > 0; i-- {
		a[i] -= a[i-1]
		c[i] -= c[i-1]
	}
	// Sort and compare
	sort.Ints(a)
	sort.Ints(c)
	for i := 0; i < n; i++ {
		if a[i] != c[i] {
			fmt.Fprint(writer, "No")
			return
		}
	}
	fmt.Fprint(writer, "Yes")
}

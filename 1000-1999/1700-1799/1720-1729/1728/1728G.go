package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement the full solution for problem G.
// Currently this program reads the input format and outputs zero for each query.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	const MOD int64 = 998244353

	var d, n, m int
	if _, err := fmt.Fscan(reader, &d, &n, &m); err != nil {
		return
	}

	lanterns := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &lanterns[i])
	}

	points := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &points[i])
	}

	var q int
	fmt.Fscan(reader, &q)
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &queries[i])
	}

	// Placeholder: output 0 for each query as the number of valid assignments.
	// Replace this with the actual combinatorial calculation.
	for i := 0; i < q; i++ {
		_ = d
		_ = n
		_ = m
		_ = lanterns
		_ = points
		_ = queries[i]
		fmt.Fprintln(writer, 0%MOD)
	}
}

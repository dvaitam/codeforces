package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Optimize I/O operations for speed (similar to ios::sync_with_stdio(0))
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var s int64

	// Read n and s
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	
	// Double the target area to work with integers (since 2 * Area = Determinant)
	s *= 2

	// Allocate slices with size n+1 to use 1-based indexing like the C++ code
	x := make([]int64, n+1)
	y := make([]int64, n+1)

	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}

	// Iterate through all unique triplets of points
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			for k := j + 1; k <= n; k++ {
				// Calculate the determinant (2 * Signed Area)
				// Formula: (y_j - y_i)(x_k - x_i) + (y_i - y_k)(x_j - x_i)
				t := (y[j]-y[i])*(x[k]-x[i]) + (y[i]-y[k])*(x[j]-x[i])

				// Check if absolute value of t equals the target 2*S
				if t == s || t == -s {
					fmt.Fprintln(writer, "Yes")
					fmt.Fprintln(writer, x[i], y[i])
					fmt.Fprintln(writer, x[j], y[j])
					fmt.Fprintln(writer, x[k], y[k])
					return
				}
			}
		}
	}

	fmt.Fprintln(writer, "No")
}

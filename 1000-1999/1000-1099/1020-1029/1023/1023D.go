package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, q int
	a    []int
	l, r []int
	fa   []int
)

// find implements the DSU find operation with path compression.
// It returns the next available index that hasn't been fully processed.
func find(i int) int {
	if fa[i] == i {
		return i
	}
	fa[i] = find(fa[i])
	return fa[i]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Use buffered I/O for efficient reading and writing
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Read N (array size) and Q (maximum value)
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	// Initialize slices
	// a: the array values (0-indexed)
	// l, r: left and right bounds for values 1..q (size q+1)
	// fa: parent array for DSU (size n+2 to handle boundary n+1)
	a = make([]int, n)
	l = make([]int, q+1)
	r = make([]int, q+1)
	fa = make([]int, n+2)

	// Initialize bounds: l[i] = n (max possible + 1), r[i] = 0 (min possible)
	for i := 0; i <= q; i++ {
		l[i] = n
		r[i] = 0
	}

	// Read the array and determine the span [l[v], r[v]] for each value v
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		val := a[i]
		l[val] = min(l[val], i)
		r[val] = max(r[val], i)
	}

	// Initialize DSU: each index points to itself initially
	for i := 0; i <= n+1; i++ {
		fa[i] = i
	}

	// Logic Step 1: Handle the case where the maximum value 'q' is missing.
	// In valid range structures, the maximum value must exist.
	// If l[q] > r[q], it means q is not in the input. We must place it.
	if l[q] > r[q] {
		// If l[0] > r[0], there are no '0' (unknown) slots to place q. Impossible.
		if l[0] > r[0] {
			fmt.Fprintln(writer, "NO")
			return
		}
		// Place q at the first available '0' slot
		idx := l[0]
		a[idx] = q
		// Mark this index as processed in DSU
		fa[idx] = find(idx + 1)
	}

	// Logic Step 2: Reverse Sweep from q down to 1.
	// We fill the ranges defined by value 'i'.
	// Any empty spot inside [l[i], r[i]] must be filled with 'i'.
	// DSU allows us to skip indices that were already filled by larger values.
	for i := q; i >= 1; i-- {
		L := l[i]
		R := r[i]

		// Iterate through indices j in [L, R] skipping already processed ones.
		// find(L) returns the first unprocessed index >= L.
		for j := find(L); j <= R; j = find(j) {
			// Check consistency:
			// If we encounter a non-zero value smaller than i inside the range of i,
			// the structure is invalid (e.g., a '1' inside a range of '2's without '2's being max).
			if a[j] != 0 && a[j] < i {
				fmt.Fprintln(writer, "NO")
				return
			}
			// Assign the value
			a[j] = i
			// Point j to j+1 to skip this index in future finds
			fa[j] = find(j + 1)
		}
	}

	// Output result
	fmt.Fprintln(writer, "YES")
	for i := 0; i < n; i++ {
		// Replace any remaining 0s with 1 (smallest valid value)
		val := a[i]
		if val == 0 {
			val = 1
		}
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, val)
	}
	writer.WriteByte('\n')
}

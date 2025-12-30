package main

import (
	"bufio"
	"fmt"
	"os"
)

// Global array for DSU (using 1<<22 as in original code)
var f [1 << 22]int

// F is the Find operation with Path Compression
func F(x int) int {
	if f[x] == 0 {
		return x
	}
	f[x] = F(f[x])
	return f[x]
}

func main() {
	// Using bufio for faster I/O in competitive programming
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)

	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)

		rootX := F(x)
		rootY := F(y + n) // Connecting row x to column y

		if rootX != rootY {
			f[rootX] = rootY
			m-- // Decrement component count
		}
	}

	// Output the final result
	fmt.Fprintln(writer, n+m-1)
}

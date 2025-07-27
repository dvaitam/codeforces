package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemC from contest 1490. For each input value x, it
// checks if x can be written as the sum of two positive cubes. Since x is at
// most 1e12, the cube root is at most 10000. We precompute all cubes up to
// 10000^3 and then test whether x-a^3 is itself a cube for some a.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	const limit = 10000
	cubes := make(map[int64]bool, limit)
	for i := 1; i <= limit; i++ {
		v := int64(i)
		cubes[v*v*v] = true
	}

	for ; t > 0; t-- {
		var x int64
		fmt.Fscan(reader, &x)
		found := false
		for i := 1; i <= limit; i++ {
			a := int64(i)
			cubeA := a * a * a
			if cubeA > x {
				break
			}
			if cubes[x-cubeA] {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

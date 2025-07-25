package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// Given the remaining sequence b, it determines if there exists
// an original sequence a of length n+k with product 2023. If so,
// it outputs k removed elements that make the product 2023.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	const target = 2023
	for ; t > 0; t-- {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return
		}
		prod := 1
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			prod *= x
		}
		if target%prod != 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		missing := target / prod
		fmt.Fprintln(writer, "YES")
		fmt.Fprint(writer, missing)
		for i := 1; i < k; i++ {
			fmt.Fprint(writer, " 1")
		}
		fmt.Fprintln(writer)
	}
}

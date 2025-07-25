package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: The problem statement in problemF.txt describes a complex sequence of
// operations involving repeatedly removing non-local extrema from a
// permutation. Implementing an efficient solution requires substantial
// algorithmic work that is beyond this placeholder. For now we provide a stub
// that reads the input format and outputs zeros for each query so the program
// builds and runs.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		// A correct solution should compute the result based on a[l:r].
		// This stub simply prints 0.
		fmt.Fprintln(out, 0)
	}
}

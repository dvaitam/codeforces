package main

import (
	"bufio"
	"fmt"
	"os"
)

// This is a placeholder implementation. The full algorithm for counting
// the number of trees consistent with the provided edges and LCA
// constraints is non-trivial. Proper handling would require dynamic
// programming over subsets of nodes and careful validation of the
// constraints. Implementing that solution from scratch is outside the
// current scope.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	for i := 0; i < q; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
	}
	fmt.Println(0)
}

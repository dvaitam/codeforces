package main

import (
	"bufio"
	"fmt"
	"os"
)

// solve is a placeholder implementation for problem G.
// The full problem statement was not available in the repository,
// so this program simply reads the input format and outputs 0 for
// each test case. Replace this with the real solution once the
// specifications are known.
func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < m; i++ {
			var v, u int
			fmt.Fscan(in, &v, &u)
		}
		fmt.Fprintln(out, 0)
	}
}

func main() {
	solve()
}

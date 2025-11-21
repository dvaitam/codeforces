package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	i, j int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		q := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &q[i])
		}

		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			pos[p[i]] = i
		}

		var ops []pair

		for i := 0; i < n; i++ {
			target := q[i]
			j := pos[target]
			for j > i {
				// swap p[j-1] and p[j]
				p[j-1], p[j] = p[j], p[j-1]
				pos[p[j]] = j
				pos[p[j-1]] = j - 1
				ops = append(ops, pair{j, j + 1}) // store 1-based positions
				j--
			}
		}

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintln(out, op.i, op.j)
		}
	}
}

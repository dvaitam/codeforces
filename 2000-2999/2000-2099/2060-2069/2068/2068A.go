package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	type pair struct{ a, b int }
	edges := make([]pair, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].a, &edges[i].b)
	}

	fmt.Fprintln(out, "YES")
	if m == 0 {
		fmt.Fprintln(out, 1)
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, i)
		}
		fmt.Fprintln(out)
		return
	}

	fmt.Fprintln(out, 2*m)

	for _, e := range edges {
		a, b := e.a, e.b
		others := make([]int, 0, n-2)
		for x := 1; x <= n; x++ {
			if x != a && x != b {
				others = append(others, x)
			}
		}

		// first vote: a, b, others ascending
		vote1 := make([]int, 0, n)
		vote1 = append(vote1, a, b)
		vote1 = append(vote1, others...)
		printPerm(out, vote1)

		// second vote: others descending, a, b
		vote2 := make([]int, 0, n)
		for i := len(others) - 1; i >= 0; i-- {
			vote2 = append(vote2, others[i])
		}
		vote2 = append(vote2, a, b)
		printPerm(out, vote2)
	}
}

func printPerm(out *bufio.Writer, perm []int) {
	for i, v := range perm {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}

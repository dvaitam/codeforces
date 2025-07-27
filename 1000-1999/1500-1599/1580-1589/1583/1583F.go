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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	// determine minimal number of colors c such that k^c >= n
	c := 0
	power := 1
	for power < n {
		power *= k
		c++
	}

	// label each node with its base-k representation of length c
	labels := make([][]int, n)
	for i := 0; i < n; i++ {
		labels[i] = make([]int, c)
		x := i
		for j := c - 1; j >= 0; j-- {
			labels[i][j] = x % k
			x /= k
		}
	}

	fmt.Fprintln(out, c)

	m := n * (n - 1) / 2
	res := make([]int, 0, m)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// find first differing digit
			color := 1
			for d := 0; d < c; d++ {
				if labels[i][d] != labels[j][d] {
					color = d + 1
					break
				}
			}
			res = append(res, color)
		}
	}

	for i := 0; i < m; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}

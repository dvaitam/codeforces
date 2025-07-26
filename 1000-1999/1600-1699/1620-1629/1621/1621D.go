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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		size := 2 * n
		// read matrix
		c := make([][]int64, size)
		for i := 0; i < size; i++ {
			c[i] = make([]int64, size)
			for j := 0; j < size; j++ {
				fmt.Fscan(in, &c[i][j])
			}
		}

		var sum int64
		for i := n; i < size; i++ {
			for j := n; j < size; j++ {
				sum += c[i][j]
			}
		}

		candidates := []int64{
			c[0][n],
			c[0][size-1],
			c[n-1][n],
			c[n-1][size-1],
			c[n][0],
			c[n][n-1],
			c[size-1][0],
			c[size-1][n-1],
		}
		best := candidates[0]
		for _, v := range candidates[1:] {
			if v < best {
				best = v
			}
		}
		fmt.Fprintln(out, sum+best)
	}
}

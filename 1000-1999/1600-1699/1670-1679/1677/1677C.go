package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ca := make([]int, n)
		cb := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &ca[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &cb[i])
		}
		fmt.Println(solve(n, ca, cb))
	}
}

func solve(n int, ca, cb []int) int {
	// build permutation mapping from color in first tape to color in second tape
	to := make([]int, n+1)
	for i := 0; i < n; i++ {
		to[ca[i]] = cb[i]
	}

	visited := make([]bool, n+1)
	pairs := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			// traverse the cycle starting from i
			cur := 0
			j := i
			for !visited[j] {
				visited[j] = true
				cur++
				j = to[j]
			}
			if cur > 1 {
				pairs += cur / 2
			}
		}
	}

	// maximum beauty is obtained by pairing extremes for each pair
	// resulting value is 2 * pairs * (n - pairs)
	return 2 * pairs * (n - pairs)
}

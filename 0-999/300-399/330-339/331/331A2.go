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

	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefixSum := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixSum[i+1] = prefixSum[i] + a[i]
	}

	maxSum := int64(-1 << 60)
	var bestL, bestR int

	for l := 0; l < n; l++ {
		for r := n - 1; r > l; r-- {
			if a[l] == a[r] {
				sum := prefixSum[r+1] - prefixSum[l]
				if sum > maxSum {
					maxSum = sum
					bestL, bestR = l, r
				}
			}
		}
	}

	if maxSum < 0 {
		maxSum = 0
		bestL, bestR = 0, 1
	}

	remove := []int{}
	for i := 0; i < bestL; i++ {
		remove = append(remove, i+1)
	}
	for i := bestR + 1; i < n; i++ {
		remove = append(remove, i+1)
	}

	fmt.Fprintf(out, "%d %d\n", maxSum, len(remove))
	for i, x := range remove {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, x)
	}
	if len(remove) > 0 {
		fmt.Fprintln(out)
	}
}

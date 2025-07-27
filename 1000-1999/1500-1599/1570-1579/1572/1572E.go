package main

import (
	"bufio"
	"fmt"
	"os"
)

// area2 computes double area of polygon from i to j inclusive (i<j).
// points are given in arrays x, y.
func area2(x, y []int64, i, j int) int64 {
	sum := int64(0)
	for t := i; t < j; t++ {
		sum += x[t]*y[t+1] - x[t+1]*y[t]
	}
	sum += x[j]*y[i] - x[i]*y[j]
	if sum < 0 {
		sum = -sum
	}
	return sum
}

// can checks if it is possible to obtain at least k+1 regions
// each having double area >= T using non-intersecting diagonals
// when restricting to contiguous sub-polygons.
func can(x, y []int64, k int, T int64) bool {
	n := len(x)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for l := 2; l < n; l++ {
		for i := 0; i+l < n; i++ {
			j := i + l
			if area2(x, y, i, j) >= T {
				dp[i][j] = 1
				for m := i + 1; m < j; m++ {
					if dp[i][m] > 0 && dp[m][j] > 0 {
						if dp[i][m]+dp[m][j] > dp[i][j] {
							dp[i][j] = dp[i][m] + dp[m][j]
						}
					}
				}
			}
		}
	}
	return dp[0][n-1] >= k+1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(reader, &n, &k)
	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}
	total := area2(x, y, 0, n-1)
	low, high := int64(0), total/int64(k+1)
	for low < high {
		mid := (low + high + 1) / 2
		if can(x, y, k, mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	fmt.Println(low)
}

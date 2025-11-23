package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Use bufio for fast I/O
	scanner := bufio.NewScanner(os.Stdin)
	
	// Read first string
	if !scanner.Scan() {
		return
	}
	s1 := scanner.Text()

	// Read second string
	if !scanner.Scan() {
		return
	}
	s2 := scanner.Text()

	n := len(s1)
	m := len(s2)

	// dp[i][j] stores the minimum cost to transform s1[:i] into s2[:j]
	// We use n+1 and m+1 to accommodate the empty prefix case (index 0)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	// --- Base Cases ---

	// 1. Transform empty string to empty string: Cost 0
	dp[0][0] = 0

	// 2. Transform s1[:i] to empty string (Deletion)
	// We must delete every character in s1 up to i
	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + charCost(s1[i-1])
	}

	// 3. Transform empty string to s2[:j] (Insertion)
	// We must insert every character in s2 up to j
	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + charCost(s2[j-1])
	}

	// --- DP Transitions ---
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			// Calculate costs for the three allowed operations:

			// 1. Substitution: replace s1[i-1] with s2[j-1]
			// Cost is absolute difference of their alphabet indices
			costSub := dp[i-1][j-1] + abs(charCost(s1[i-1])-charCost(s2[j-1]))

			// 2. Deletion: delete s1[i-1]
			// We move from state [i-1][j] to [i][j]
			costDel := dp[i-1][j] + charCost(s1[i-1])

			// 3. Insertion: insert s2[j-1]
			// We move from state [i][j-1] to [i][j]
			costIns := dp[i][j-1] + charCost(s2[j-1])

			// Take the minimum of the three
			dp[i][j] = min(costSub, costDel, costIns)
		}
	}

	// The answer is in the bottom-right cell
	fmt.Println(dp[n][m])
}

// charCost returns the 1-based index of a lowercase latin letter (a=1, b=2...)
func charCost(b byte) int {
	return int(b - 'a' + 1)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, p, k int
	if _, err := fmt.Fscan(in, &n, &p, &k); err != nil {
		return
	}
	X := make([]int, n)
	Y := make([]int, n)
	var r int
	fmt.Fscan(in, &r)
	for i := 0; i < r; i++ {
		var a int
		fmt.Fscan(in, &a)
		if a >= 1 && a <= n {
			X[a-1] = 1
		}
	}
	fmt.Fscan(in, &r)
	for i := 0; i < r; i++ {
		var b int
		fmt.Fscan(in, &b)
		if b >= 1 && b <= n {
			Y[b-1] = 1
		}
	}

	if 2*((n+k-1)/k) < p {
		p = 2 * ((n + k - 1) / k)
	}
	if p == 0 {
		// cannot peek at all
		fmt.Println(0)
		return
	}

	const NEG = -1 << 60
	// dp[2][p+1][k][k]
	dp := make([][][][]int, 2)
	for t := 0; t < 2; t++ {
		dp[t] = make([][][]int, p+1)
		for j := 0; j <= p; j++ {
			dp[t][j] = make([][]int, k)
			for a := 0; a < k; a++ {
				dp[t][j][a] = make([]int, k)
				for b := 0; b < k; b++ {
					dp[t][j][a][b] = NEG
				}
			}
		}
	}

	dp[0][0][0][0] = 0

	for i := 1; i <= n; i++ {
		curr := i % 2
		prev := 1 - curr
		// reset current layer
		for j := 0; j <= p; j++ {
			for a := 0; a < k; a++ {
				for b := 0; b < k; b++ {
					dp[curr][j][a][b] = NEG
				}
			}
		}
		xi := X[i-1]
		yi := Y[i-1]
		// use overflow from previous
		for j := 1; j <= p; j++ {
			for c := 0; c < k-1; c++ {
				if dp[prev][j][c+1][0] != NEG {
					val := dp[prev][j][c+1][0] + xi
					if val > dp[curr][j][c][0] {
						dp[curr][j][c][0] = val
					}
				}
				if dp[prev][j][0][c+1] != NEG {
					val := dp[prev][j][0][c+1] + yi
					if val > dp[curr][j][0][c] {
						dp[curr][j][0][c] = val
					}
				}
			}
		}
		// start peek at i
		for j := 1; j <= p; j++ {
			if dp[prev][j-1][0][0] != NEG {
				val := dp[prev][j-1][0][0] + xi
				if val > dp[curr][j][k-1][0] {
					dp[curr][j][k-1][0] = val
				}
				val2 := dp[prev][j-1][0][0] + yi
				if val2 > dp[curr][j][0][k-1] {
					dp[curr][j][0][k-1] = val2
				}
			}
		}
		// overflow from both
		for j := 2; j <= p; j++ {
			for a := 0; a < k-1; a++ {
				for b := 0; b < k-1; b++ {
					if dp[prev][j][a+1][b+1] != NEG {
						val := dp[prev][j][a+1][b+1] + boolToInt(xi == 1 || yi == 1)
						if val > dp[curr][j][a][b] {
							dp[curr][j][a][b] = val
						}
					}
				}
			}
		}
		// one overflow, one new peek
		for j := 1; j <= p; j++ {
			for c := 0; c < k-1; c++ {
				if dp[prev][j-1][0][c+1] != NEG {
					val := dp[prev][j-1][0][c+1] + boolToInt(xi == 1 || yi == 1)
					if val > dp[curr][j][k-1][c] {
						dp[curr][j][k-1][c] = val
					}
				}
				if dp[prev][j-1][c+1][0] != NEG {
					val := dp[prev][j-1][c+1][0] + boolToInt(xi == 1 || yi == 1)
					if val > dp[curr][j][c][k-1] {
						dp[curr][j][c][k-1] = val
					}
				}
			}
		}
		// no peek, just carry over
		for j := 0; j <= p; j++ {
			if dp[prev][j][0][0] > dp[curr][j][0][0] {
				dp[curr][j][0][0] = dp[prev][j][0][0]
			}
		}
		// clear prev layer
		for j := 0; j <= p; j++ {
			for a := 0; a < k; a++ {
				for b := 0; b < k; b++ {
					dp[prev][j][a][b] = NEG
				}
			}
		}
	}

	ans := 0
	last := n % 2
	for j := 1; j <= p; j++ {
		for a := 0; a < k; a++ {
			for b := 0; b < k; b++ {
				if dp[last][j][a][b] > ans {
					ans = dp[last][j][a][b]
				}
			}
		}
	}
	fmt.Println(ans)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

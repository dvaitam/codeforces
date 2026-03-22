package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	t := make([]int64, n+2)
	x := make([]int64, n+2)

	t[0] = 0
	x[0] = 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &t[i], &x[i])
	}
	const INF int64 = 2e18
	t[n+1] = INF
	x[n+1] = 0

	can_seq := make([]bool, n+2)
	for k := 1; k < n; k++ {
		can_seq[k] = (t[k]+abs(x[k]-x[k+1]) <= t[k+1])
	}

	pos := make([]bool, n+2)
	dp := make([]int64, n+2)
	for i := 0; i <= n+1; i++ {
		dp[i] = INF
	}
	pos[0] = true

	ans := false

	for i := 0; i < n; i++ {
		if pos[i] {
			if t[i]+abs(x[i]-x[i+1]) <= t[i+1] {
				if i+1 == n {
					ans = true
				} else {
					pos[i+1] = true
					val := t[i] + abs(x[i]-x[i+1]) + abs(x[i+1]-x[i+2])
					if val < dp[i+1] {
						dp[i+1] = val
					}
				}
			}

			isValid := true
			for j := i + 2; j <= n; j++ {
				if j-2 >= i+1 {
					if !can_seq[j-2] {
						isValid = false
					}
				}
				if !isValid {
					break
				}

				T := t[i] + abs(x[i]-x[j]) + abs(x[j]-x[i+1])
				if T <= t[i+1] {
					if j == n {
						ans = true
					} else {
						val := t[j-1] + abs(x[j-1]-x[j+1])
						if val < dp[j] {
							dp[j] = val
						}
					}
				}
			}
		}

		if dp[i] != INF {
			if dp[i] <= t[i+1] {
				if i+1 == n {
					ans = true
				} else {
					pos[i+1] = true
					val := max(t[i], dp[i]) + abs(x[i+1]-x[i+2])
					if val < dp[i+1] {
						dp[i+1] = val
					}
				}
			}

			isValid := true
			for j := i + 2; j <= n; j++ {
				if j-2 >= i+1 {
					if !can_seq[j-2] {
						isValid = false
					}
				}
				if !isValid {
					break
				}

				T := max(t[i], dp[i]+abs(x[i+1]-x[j])) + abs(x[j]-x[i+1])
				if T <= t[i+1] {
					if j == n {
						ans = true
					} else {
						val := t[j-1] + abs(x[j-1]-x[j+1])
						if val < dp[j] {
							dp[j] = val
						}
					}
				}
			}
		}
	}

	if ans {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	oddExisting := 0
	evenExisting := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		if arr[i] != 0 {
			if arr[i]%2 == 1 {
				oddExisting++
			} else {
				evenExisting++
			}
		}
	}

	totOdd := (n + 1) / 2
	totEven := n / 2
	oddMissing := totOdd - oddExisting
	evenMissing := totEven - evenExisting

	// prefix zeros count
	prefixZero := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefixZero[i] = prefixZero[i-1]
		if arr[i-1] == 0 {
			prefixZero[i]++
		}
	}

	const INF int = 1 << 30
	dp := make([][][]int, n+1)
	for i := range dp {
		dp[i] = make([][]int, oddMissing+1)
		for j := range dp[i] {
			dp[i][j] = []int{INF, INF}
		}
	}
	dp[0][0][0] = 0
	dp[0][0][1] = 0

	for i := 0; i < n; i++ {
		for o := 0; o <= oddMissing; o++ {
			for last := 0; last < 2; last++ {
				cur := dp[i][o][last]
				if cur == INF {
					continue
				}
				zeroSoFar := prefixZero[i]
				evenUsed := zeroSoFar - o
				if arr[i] != 0 {
					p := arr[i] % 2
					add := 0
					if i > 0 && last != p {
						add = 1
					}
					if dp[i+1][o][p] > cur+add {
						dp[i+1][o][p] = cur + add
					}
				} else {
					// place odd number
					if o < oddMissing {
						add := 0
						if i > 0 && last != 1 {
							add = 1
						}
						if dp[i+1][o+1][1] > cur+add {
							dp[i+1][o+1][1] = cur + add
						}
					}
					// place even number
					if evenUsed < evenMissing {
						add := 0
						if i > 0 && last != 0 {
							add = 1
						}
						if dp[i+1][o][0] > cur+add {
							dp[i+1][o][0] = cur + add
						}
					}
				}
			}
		}
	}

	res := dp[n][oddMissing][0]
	if dp[n][oddMissing][1] < res {
		res = dp[n][oddMissing][1]
	}
	fmt.Fprintln(writer, res)
}

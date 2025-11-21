package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	children := make([][]int, n)
	value := make([]int64, n)
	root := 0

	for i := 0; i < n; i++ {
		var p int
		var a int64
		fmt.Fscan(in, &p, &a)
		value[i] = a
		if p == -1 {
			root = i
		} else {
			children[p-1] = append(children[p-1], i)
		}
	}

	dp := make([][2]int64, n)

	for i := n - 1; i >= 0; i-- {
		arr := [2]int64{0, negInf}
		for _, ch := range children[i] {
			next := [2]int64{negInf, negInf}
			for parity := 0; parity < 2; parity++ {
				if arr[parity] == negInf {
					continue
				}
				for childParity := 0; childParity < 2; childParity++ {
					val := dp[ch][childParity]
					if val == negInf {
						continue
					}
					np := parity ^ childParity
					if cand := arr[parity] + val; cand > next[np] {
						next[np] = cand
					}
				}
			}
			arr = next
		}

		dp[i][0] = arr[0]
		bestOdd := arr[1]
		if arr[0] != negInf {
			if cand := arr[0] + value[i]; cand > bestOdd {
				bestOdd = cand
			}
		}
		dp[i][1] = bestOdd
	}

	ans := dp[root][0]
	if dp[root][1] > ans {
		ans = dp[root][1]
	}

	fmt.Fprintln(out, ans)
}

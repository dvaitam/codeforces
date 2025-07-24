package main

import (
	"bufio"
	"fmt"
	"os"
)

// countPairs returns the number of pairs (a,b) such that
// 0 <= a <= x, 0 <= b <= y and (a & b) == 0.
func countPairs(x, y int64) int64 {
	if x < 0 || y < 0 {
		return 0
	}
	var dp [2][2]int64
	dp[1][1] = 1
	for bit := 30; bit >= 0; bit-- {
		var ndp [2][2]int64
		bx := int((x >> bit) & 1)
		by := int((y >> bit) & 1)
		for tx := 0; tx <= 1; tx++ {
			for ty := 0; ty <= 1; ty++ {
				val := dp[tx][ty]
				if val == 0 {
					continue
				}
				for a := 0; a <= 1; a++ {
					if tx == 1 && a > bx {
						continue
					}
					ntx := 0
					if tx == 1 && a == bx {
						ntx = 1
					}
					for b := 0; b <= 1; b++ {
						if a == 1 && b == 1 {
							continue
						}
						if ty == 1 && b > by {
							continue
						}
						nty := 0
						if ty == 1 && b == by {
							nty = 1
						}
						ndp[ntx][nty] += val
					}
				}
			}
		}
		dp = ndp
	}
	return dp[0][0] + dp[0][1] + dp[1][0] + dp[1][1]
}

func solve(l, r int64) int64 {
	return countPairs(r, r) - countPairs(l-1, r) - countPairs(r, l-1) + countPairs(l-1, l-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		fmt.Fprintln(writer, solve(l, r))
	}
}

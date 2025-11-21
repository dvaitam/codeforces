package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxBit   = 61
	bitLimit = 59
	limitK   = 1_000_000_000_000_000_000
	inf      = int64(^uint64(0)>>1) / 4
)

func solve(x, y int64) int64 {
	if x == y {
		return -1
	}

	var dp [maxBit + 1][2][2]int64
	for i := 0; i <= maxBit; i++ {
		for cx := 0; cx < 2; cx++ {
			for cy := 0; cy < 2; cy++ {
				dp[i][cx][cy] = inf
			}
		}
	}
	dp[0][0][0] = 0

	for i := 0; i < maxBit; i++ {
		xBit := (x >> uint(i)) & 1
		yBit := (y >> uint(i)) & 1
		for cx := 0; cx < 2; cx++ {
			for cy := 0; cy < 2; cy++ {
				cur := dp[i][cx][cy]
				if cur == inf {
					continue
				}

				maxKb := 0
				if i <= bitLimit {
					maxKb = 1
				}
				for kb := 0; kb <= maxKb; kb++ {
					totalX := int(xBit) + kb + cx
					nextCx := totalX >> 1
					sumX := totalX & 1

					totalY := int(yBit) + kb + cy
					nextCy := totalY >> 1
					sumY := totalY & 1

					if sumX == 1 && sumY == 1 {
						continue
					}

					nextVal := cur + (int64(kb) << uint(i))
					if nextVal < dp[i+1][nextCx][nextCy] {
						dp[i+1][nextCx][nextCy] = nextVal
					}
				}
			}
		}
	}

	ans := dp[maxBit][0][0]
	if ans == inf || ans > limitK {
		return -1
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var x, y int64
		fmt.Fscan(in, &x, &y)
		fmt.Fprintln(out, solve(x, y))
	}
}

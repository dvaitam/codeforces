package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func modPow2(exp int64) int64 {
	res := int64(1)
	base := int64(2)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		rowCols := make([][]int, n)
		totalStars := int64(0)
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				if grid[r][c] == '1' {
					rowCols[r] = append(rowCols[r], c)
					totalStars++
				}
			}
		}

		if totalStars == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		lim := 1 << n // n <= 14
		dpCnt := make([]int64, lim)
		dpSum := make([]int64, lim)
		dpCnt[0] = 1

		for r := 0; r < n; r++ {
			if len(rowCols[r]) == 0 {
				continue
			}
			newCnt := make([]int64, lim)
			newSum := make([]int64, lim)
			copy(newCnt, dpCnt)
			copy(newSum, dpSum)
			for mask := 0; mask < lim; mask++ {
				if dpCnt[mask] == 0 {
					continue
				}
				for _, c := range rowCols[r] {
					if mask&(1<<c) != 0 {
						continue
					}
					nm := mask | (1 << c)
					newCnt[nm] = (newCnt[nm] + dpCnt[mask]) % mod
					newSum[nm] = (newSum[nm] + dpSum[mask] + dpCnt[mask]) % mod // size +1
				}
			}
			dpCnt, newCnt = newCnt, dpCnt
			dpSum, newSum = newSum, dpSum
		}

		var matchCnt, matchSizeSum int64
		for i := 0; i < lim; i++ {
			matchCnt = (matchCnt + dpCnt[i]) % mod
			matchSizeSum = (matchSizeSum + dpSum[i]) % mod
		}

		totalSubsets := modPow2(totalStars)
		ans := (totalSubsets - matchCnt + matchSizeSum) % mod
		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(out, ans)
	}
}

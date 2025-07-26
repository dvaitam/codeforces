package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func solveCase(n, k int, colors []int) int {
	freq := make([]int, n+1)
	for _, c := range colors {
		freq[c]++
	}
	lens := make([][]int, n+1)
	cnts := make([][]int, n+1)
	for c := 1; c <= n; c++ {
		if freq[c] > 0 {
			size := k
			if freq[c] < size {
				size = freq[c]
			}
			lens[c] = make([]int, size+1)
			cnts[c] = make([]int, size+1)
			for i := 1; i <= size; i++ {
				lens[c][i] = -1
			}
		}
	}
	bestLen := 0
	bestCnt := 1
	for _, col := range colors {
		arrLen := lens[col]
		arrCnt := cnts[col]
		size := len(arrLen) - 1
		oldKLen, oldKCnt := 0, 0
		if k <= size {
			oldKLen = arrLen[k]
			oldKCnt = arrCnt[k]
		}
		for r := size; r >= 1; r-- {
			baseLen := -1
			baseCnt := 0
			if r == 1 {
				baseLen = bestLen
				baseCnt = bestCnt
			} else if arrLen[r-1] != -1 {
				baseLen = arrLen[r-1]
				baseCnt = arrCnt[r-1]
			}
			if baseLen != -1 {
				newLen := baseLen + 1
				if newLen > arrLen[r] {
					arrLen[r] = newLen
					arrCnt[r] = baseCnt
				} else if newLen == arrLen[r] {
					arrCnt[r] += baseCnt
					if arrCnt[r] >= MOD {
						arrCnt[r] -= MOD
					}
				}
			}
		}
		if k <= size {
			if arrLen[k] > bestLen {
				bestLen = arrLen[k]
				bestCnt = arrCnt[k] % MOD
			} else if arrLen[k] == bestLen {
				added := arrCnt[k]
				if oldKLen == bestLen {
					added -= oldKCnt
					added %= MOD
					if added < 0 {
						added += MOD
					}
				}
				bestCnt += added
				bestCnt %= MOD
			}
		}
	}
	return bestCnt % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &colors[i])
		}
		ans := solveCase(n, k, colors)
		fmt.Fprintln(writer, ans)
	}
}

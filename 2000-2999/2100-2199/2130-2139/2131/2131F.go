package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var aStr, bStr string
		fmt.Fscan(in, &aStr)
		fmt.Fscan(in, &bStr)

		prefA0 := make([]int64, n+1)
		prefA1 := make([]int64, n+1)
		prefB0 := make([]int64, n+1)
		prefB1 := make([]int64, n+1)

		for i := 1; i <= n; i++ {
			if aStr[i-1] == '0' {
				prefA0[i] = prefA0[i-1] + 1
				prefA1[i] = prefA1[i-1]
			} else {
				prefA0[i] = prefA0[i-1]
				prefA1[i] = prefA1[i-1] + 1
			}
			if bStr[i-1] == '0' {
				prefB0[i] = prefB0[i-1] + 1
				prefB1[i] = prefB1[i-1]
			} else {
				prefB0[i] = prefB0[i-1]
				prefB1[i] = prefB1[i-1] + 1
			}
		}

		sumA0 := int64(0)
		sumB0 := int64(0)
		deltaA := make([]int64, n+1)
		deltaB := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			sumA0 += prefA0[i]
			sumB0 += prefB0[i]
			deltaA[i] = prefA1[i] - prefA0[i]
			deltaB[i] = prefB1[i] - prefB0[i]
		}

		size := 2*n + 5
		offset := n + 2
		freq := make([]int64, size)
		for i := 1; i <= n; i++ {
			freq[int(deltaB[i])+offset]++
		}
		prefCnt := make([]int64, size)
		prefSum := make([]int64, size)
		prefCnt[0] = freq[0]
		prefSum[0] = freq[0] * int64(0-offset)
		for i := 1; i < size; i++ {
			prefCnt[i] = prefCnt[i-1] + freq[i]
			prefSum[i] = prefSum[i-1] + freq[i]*int64(i-offset)
		}

		query := func(th int64) (int64, int64) {
			idx := int(th) + offset
			if idx < 0 {
				return 0, 0
			}
			if idx >= size {
				return prefCnt[size-1], prefSum[size-1]
			}
			return prefCnt[idx], prefSum[idx]
		}

		total := int64(n)*sumA0 + int64(n)*sumB0
		for i := 1; i <= n; i++ {
			cnt, sumVal := query(-deltaA[i])
			total += cnt*deltaA[i] + sumVal
		}

		fmt.Fprintln(out, total)
	}
}

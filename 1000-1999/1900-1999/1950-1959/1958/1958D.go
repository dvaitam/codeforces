package main

import (
	"bufio"
	"fmt"
	"os"
)

func costBlock(arr []int64) int64 {
	m := len(arr)
	if m%2 == 0 {
		var cost int64
		for i := 0; i < m; i += 2 {
			cost += 2 * (arr[i] + arr[i+1])
		}
		return cost
	}

	pairOdd := make([]int64, m+1)
	pairEven := make([]int64, m+1)
	for i := 1; i < m; i++ {
		pairOdd[i] = pairOdd[i-1]
		pairEven[i] = pairEven[i-1]
		if i%2 == 1 {
			pairOdd[i] += 2 * (arr[i-1] + arr[i])
		} else {
			pairEven[i] += 2 * (arr[i-1] + arr[i])
		}
	}
	pairOdd[m] = pairOdd[m-1]
	pairEven[m] = pairEven[m-1]

	ans := int64(1<<62 - 1)
	for j := 1; j <= m; j += 2 {
		cost := pairOdd[j-1] + arr[j-1] + (pairEven[m] - pairEven[j])
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		var total int64
		for i := 0; i < n; {
			if a[i] == 0 {
				i++
				continue
			}
			start := i
			for i < n && a[i] > 0 {
				i++
			}
			total += costBlock(a[start:i])
		}
		fmt.Fprintln(writer, total)
	}
}

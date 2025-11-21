package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m, q int
		fmt.Fscan(reader, &n, &m, &q)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

		prefA := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefA[i] = prefA[i-1] + a[i-1]
		}
		prefB := make([]int64, m+1)
		for i := 1; i <= m; i++ {
			prefB[i] = prefB[i-1] + b[i-1]
		}

		for ; q > 0; q-- {
			var x, y, z int
			fmt.Fscan(reader, &x, &y, &z)
			ans := int64(0)
			maxA := x
			if maxA > n {
				maxA = n
			}
			maxB := y
			if maxB > m {
				maxB = m
			}
			if z > maxA+maxB {
				z = maxA + maxB
			}
			for takeA := 0; takeA <= z && takeA <= maxA; takeA++ {
				takeB := z - takeA
				if takeB < 0 || takeB > maxB {
					continue
				}
				current := prefA[takeA] + prefB[takeB]
				if current > ans {
					ans = current
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}

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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &a[i])
		}
		L := 1
		if n-L < 1 {
			L = n - (n - 1)
		}
		var ans int64
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				if i == j {
					continue
				}
				needMin := n - a[j]
				if needMin < 1 {
					needMin = 1
				}
				maxVal := a[i]
				if maxVal > n-1 {
					maxVal = n - 1
				}
				if needMin <= maxVal {
					ans += int64(maxVal - needMin + 1)
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}

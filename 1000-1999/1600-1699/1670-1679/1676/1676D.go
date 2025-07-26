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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &a[i][j])
			}
		}
		d1 := make([]int64, n+m)
		d2 := make([]int64, n+m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				d1[i+j] += a[i][j]
				d2[i-j+m-1] += a[i][j]
			}
		}
		var ans int64
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				sum := d1[i+j] + d2[i-j+m-1] - a[i][j]
				if sum > ans {
					ans = sum
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

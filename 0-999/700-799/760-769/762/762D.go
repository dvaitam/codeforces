package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([][]int, 3)
	for i := 0; i < 3; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	dp0 := a[0][0]
	dp1 := a[0][0] + a[1][0]
	dp2 := a[0][0] + a[1][0] + a[2][0]
	for i := 1; i < n; i++ {
		ndp0 := max(dp0+a[0][i], max(dp1+a[0][i]+a[1][i], dp2+a[0][i]+a[1][i]+a[2][i]))
		ndp1 := max(dp1+a[1][i], max(dp0+a[0][i]+a[1][i], dp2+a[1][i]+a[2][i]))
		ndp2 := max(dp2+a[2][i], max(dp1+a[1][i]+a[2][i], dp0+a[0][i]+a[1][i]+a[2][i]))
		dp0, dp1, dp2 = ndp0, ndp1, ndp2
	}
	fmt.Println(dp2)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

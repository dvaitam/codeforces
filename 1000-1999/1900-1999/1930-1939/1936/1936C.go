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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
		// naive approach: directly hire pokemon n to beat current champion
		diff := int64(1<<62 - 1)
		for j := 0; j < m; j++ {
			d := a[0][j] - a[n-1][j]
			if d < 0 {
				d = 0
			}
			if d < diff {
				diff = d
			}
		}
		ans := diff + c[n-1]
		fmt.Fprintln(out, ans)
	}
}

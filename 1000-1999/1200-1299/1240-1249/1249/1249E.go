package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var c int64
	if _, err := fmt.Fscan(in, &n, &c); err != nil {
		return
	}
	a := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &b[i])
	}

	stairs := make([]int64, n)
	elev := make([]int64, n)
	for i := range stairs {
		stairs[i] = 1 << 60
		elev[i] = 1 << 60
	}
	stairs[0] = 0
	elev[0] = c
	ans := make([]int64, n)
	ans[0] = 0
	for i := 1; i < n; i++ {
		stairs[i] = min(stairs[i-1], elev[i-1]) + a[i-1]
		elev[i] = min(elev[i-1]+b[i-1], stairs[i-1]+c+b[i-1])
		ans[i] = min(stairs[i], elev[i])
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}

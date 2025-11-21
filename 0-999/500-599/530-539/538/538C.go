package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var m int
	fmt.Fscan(in, &n, &m)

	d := make([]int64, m)
	h := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &d[i], &h[i])
	}

	possible := true
	ans := int64(0)

	if m > 0 {
		ans = max64(ans, h[0]+(d[0]-1))
		ans = max64(ans, h[m-1]+(n-d[m-1]))
	}

	for i := 0; i < m-1; i++ {
		diff := abs(h[i+1] - h[i])
		dist := d[i+1] - d[i]
		if diff > dist {
			possible = false
			break
		}
		remain := dist - diff
		peak := max64(h[i], h[i+1]) + (remain / 2)
		ans = max64(ans, peak)
	}

	if !possible {
		fmt.Fprintln(out, "IMPOSSIBLE")
		return
	}

	fmt.Fprintln(out, ans)
}

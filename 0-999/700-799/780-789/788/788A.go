package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	if n < 2 {
		fmt.Fprintln(writer, 0)
		return
	}

	diff := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		diff[i] = abs64(a[i+1] - a[i])
	}

	dpPlus := make([]int64, n-1)
	dpMinus := make([]int64, n-1)

	dpPlus[0] = diff[0]
	dpMinus[0] = -diff[0]
	ans := max64(dpPlus[0], dpMinus[0])

	for i := 1; i < n-1; i++ {
		dpPlus[i] = max64(diff[i], dpMinus[i-1]+diff[i])
		dpMinus[i] = max64(-diff[i], dpPlus[i-1]-diff[i])
		ans = max64(ans, max64(dpPlus[i], dpMinus[i]))
	}

	fmt.Fprintln(writer, ans)
}

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

func solve(a []int64) int64 {
	arr := make([]int64, len(a)+2)
	copy(arr[1:], a)
	// arr[0] and arr[len(a)+1] are zero by default
	ans := int64(0)
	for i := 1; i < len(arr); i++ {
		ans += abs64(arr[i] - arr[i-1])
	}
	for i := 1; i < len(arr)-1; i++ {
		drop := arr[i] - max64(arr[i-1], arr[i+1])
		if drop > 0 {
			ans -= drop
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t++ {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, solve(arr))
	}
}

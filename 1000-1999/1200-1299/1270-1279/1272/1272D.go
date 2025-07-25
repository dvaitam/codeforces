package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	if n == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	left := make([]int, n)
	left[0] = 1
	for i := 1; i < n; i++ {
		if a[i] > a[i-1] {
			left[i] = left[i-1] + 1
		} else {
			left[i] = 1
		}
	}
	right := make([]int, n)
	right[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		if a[i] < a[i+1] {
			right[i] = right[i+1] + 1
		} else {
			right[i] = 1
		}
	}
	ans := 1
	for i := 0; i < n; i++ {
		ans = max(ans, left[i])
	}
	if n > 1 {
		ans = max(ans, right[1])
		ans = max(ans, left[n-2])
	}
	for i := 1; i+1 < n; i++ {
		if a[i+1] > a[i-1] {
			ans = max(ans, left[i-1]+right[i+1])
		}
	}
	if ans > n {
		ans = n
	}
	fmt.Fprintln(writer, ans)
}

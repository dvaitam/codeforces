package main

import (
	"bufio"
	"fmt"
	"os"
)

func canForm(a []int) bool {
	n := len(a)
	var dfs func(int, int, bool) bool
	dfs = func(i, sum int, used bool) bool {
		if i == n {
			return used && sum == 0
		}
		if dfs(i+1, sum, used) {
			return true
		}
		if dfs(i+1, sum+a[i], true) {
			return true
		}
		if dfs(i+1, sum-a[i], true) {
			return true
		}
		return false
	}
	return dfs(0, 0, false)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if canForm(arr) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

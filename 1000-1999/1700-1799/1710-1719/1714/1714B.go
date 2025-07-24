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
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		seen := make(map[int]bool)
		ans := 0
		for i := n - 1; i >= 0; i-- {
			if seen[arr[i]] {
				ans = i + 1
				break
			}
			seen[arr[i]] = true
		}
		fmt.Fprintln(writer, ans)
	}
}

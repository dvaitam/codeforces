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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(reader, &n, &x)
		arr := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			sum += arr[i]
		}
		if sum%x != 0 {
			fmt.Fprintln(writer, n)
			continue
		}
		left := -1
		for i := 0; i < n; i++ {
			if arr[i]%x != 0 {
				left = i
				break
			}
		}
		if left == -1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		right := -1
		for i := n - 1; i >= 0; i-- {
			if arr[i]%x != 0 {
				right = i
				break
			}
		}
		removeLen := left + 1
		if n-right < removeLen {
			removeLen = n - right
		}
		fmt.Fprintln(writer, n-removeLen)
	}
}

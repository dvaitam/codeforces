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
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			sum += arr[i]
		}
		possible := false
		for i := 0; i < n; i++ {
			if arr[i]*n == sum {
				possible = true
				break
			}
		}
		if possible {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

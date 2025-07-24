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
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		// compress consecutive duplicates
		b := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if i == 0 || arr[i] != arr[i-1] {
				b = append(b, arr[i])
			}
		}
		m := len(b)
		idx := 0
		for idx+1 < m && b[idx] > b[idx+1] {
			idx++
		}
		for idx+1 < m && b[idx] < b[idx+1] {
			idx++
		}
		if idx == m-1 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

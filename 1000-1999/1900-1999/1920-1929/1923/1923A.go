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
		first, last := -1, -1
		for i := 0; i < n; i++ {
			if arr[i] == 1 {
				if first == -1 {
					first = i
				}
				last = i
			}
		}
		if first == -1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		count := 0
		for i := first; i <= last; i++ {
			if arr[i] == 0 {
				count++
			}
		}
		fmt.Fprintln(writer, count)
	}
}

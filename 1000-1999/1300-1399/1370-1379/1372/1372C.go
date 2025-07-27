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

		sorted := true
		for i := 0; i < n; i++ {
			if arr[i] != i+1 {
				sorted = false
				break
			}
		}
		if sorted {
			fmt.Fprintln(writer, 0)
			continue
		}

		l := 0
		for l < n && arr[l] == l+1 {
			l++
		}
		r := n - 1
		for r >= 0 && arr[r] == r+1 {
			r--
		}

		needTwo := false
		for i := l; i <= r; i++ {
			if arr[i] == i+1 {
				needTwo = true
				break
			}
		}
		if needTwo {
			fmt.Fprintln(writer, 2)
		} else {
			fmt.Fprintln(writer, 1)
		}
	}
}

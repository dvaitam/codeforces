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

	var n, q int
	fmt.Fscan(reader, &n, &q)

	arr := make([]int, n)
	ones := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		if arr[i] == 1 {
			ones++
		}
	}

	for i := 0; i < q; i++ {
		var t, x int
		fmt.Fscan(reader, &t, &x)
		if t == 1 {
			idx := x - 1
			if arr[idx] == 1 {
				arr[idx] = 0
				ones--
			} else {
				arr[idx] = 1
				ones++
			}
		} else if t == 2 {
			if x <= ones {
				fmt.Fprintln(writer, 1)
			} else {
				fmt.Fprintln(writer, 0)
			}
		}
	}
}

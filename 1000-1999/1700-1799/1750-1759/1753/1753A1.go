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
		if n%2 == 1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		segments := make([][2]int, 0, n)
		for i := 0; i < n; i += 2 {
			if arr[i] == arr[i+1] {
				segments = append(segments, [2]int{i + 1, i + 2})
			} else {
				segments = append(segments, [2]int{i + 1, i + 1})
				segments = append(segments, [2]int{i + 2, i + 2})
			}
		}
		fmt.Fprintln(writer, len(segments))
		for _, seg := range segments {
			fmt.Fprintln(writer, seg[0], seg[1])
		}
	}
}

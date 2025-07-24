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
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var ans int64
	for l := 0; l < n; l++ {
		minv := arr[l]
		maxv := arr[l]
		for r := l; r < n; r++ {
			if arr[r] < minv {
				minv = arr[r]
			}
			if arr[r] > maxv {
				maxv = arr[r]
			}
			if maxv%minv == 0 {
				ans++
			}
		}
	}
	fmt.Fprintln(writer, ans)
}

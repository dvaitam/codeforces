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
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	minVal := arr[0]
	for i := 1; i < n; i++ {
		if arr[i] < minVal {
			minVal = arr[i]
		}
	}
	for i := 0; i < q; i++ {
		var k int64
		fmt.Fscan(reader, &k)
		// TODO: implement actual algorithm
		fmt.Fprintln(writer, minVal)
	}
}

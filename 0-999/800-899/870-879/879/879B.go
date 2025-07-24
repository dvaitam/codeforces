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

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	current := arr[0]
	var wins int64 = 0

	for i := 1; i < n && wins < k; i++ {
		if arr[i] > current {
			current = arr[i]
			wins = 1
		} else {
			wins++
		}
	}

	fmt.Fprintln(writer, current)
}

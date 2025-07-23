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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	size := 2 * n
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	swaps := 0
	for i := 0; i < size; i += 2 {
		if arr[i] == arr[i+1] {
			continue
		}
		// find partner for arr[i]
		j := i + 1
		for j < size && arr[j] != arr[i] {
			j++
		}
		// bring arr[j] next to arr[i]
		for j > i+1 {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			j--
			swaps++
		}
	}

	fmt.Fprintln(writer, swaps)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		arr := make([]int, 2*n)
		for i := range arr {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Ints(arr)
		sum := 0
		for i := 0; i < 2*n; i += 2 {
			sum += arr[i]
		}
		fmt.Fprintln(writer, sum)
	}
}

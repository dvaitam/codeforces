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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		prefix := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			sum += int64(arr[i])
			prefix[i] = sum
		}
		for j := 0; j < q; j++ {
			var x int64
			fmt.Fscan(reader, &x)
			idx := sort.Search(len(prefix), func(i int) bool {
				return prefix[i] >= x
			})
			if idx == len(prefix) {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, idx+1)
			}
		}
	}
}

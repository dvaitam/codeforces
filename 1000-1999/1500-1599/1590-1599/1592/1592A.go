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
		var n int
		var H int64
		fmt.Fscan(reader, &n, &H)
		arr := make([]int64, n)
		for i := range arr {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i] > arr[j]
		})
		a, b := arr[0], arr[1]
		pair := a + b
		cnt := H / pair * 2
		rem := H % pair
		if rem > 0 {
			if rem <= a {
				cnt++
			} else {
				cnt += 2
			}
		}
		fmt.Fprintln(writer, cnt)
	}
}

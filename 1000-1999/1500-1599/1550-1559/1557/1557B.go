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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		pos := make(map[int]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			pos[arr[i]] = i
		}
		sortedArr := make([]int, n)
		copy(sortedArr, arr)
		sort.Ints(sortedArr)
		segments := 1
		for i := 0; i < n-1; i++ {
			if pos[sortedArr[i]]+1 != pos[sortedArr[i+1]] {
				segments++
			}
		}
		if segments <= k {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

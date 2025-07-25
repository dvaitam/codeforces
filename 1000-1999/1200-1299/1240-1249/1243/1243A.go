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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Sort(sort.Reverse(sort.IntSlice(arr)))
		ans := 0
		for i := 0; i < n; i++ {
			if arr[i] >= i+1 {
				ans = i + 1
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

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
		var l, r int64
		fmt.Fscan(reader, &n, &l, &r)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

		countPairs := func(limit int64) int64 {
			var cnt int64
			j := n - 1
			for i := 0; i < n; i++ {
				for j > i && arr[i]+arr[j] > limit {
					j--
				}
				if j <= i {
					break
				}
				cnt += int64(j - i)
			}
			return cnt
		}

		res := countPairs(r) - countPairs(l-1)
		fmt.Fprintln(writer, res)
	}
}

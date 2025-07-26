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
		xs := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &xs[i])
		}

		type pair struct {
			val int64
			idx int
		}
		arr := make([]pair, n)
		for i, v := range xs {
			arr[i] = pair{v, i}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + arr[i].val
		}
		ans := make([]int64, n)
		total := int64(n)
		for i := 0; i < n; i++ {
			s := arr[i].val
			left := int64(i)*s - prefix[i]
			right := (prefix[n] - prefix[i+1]) - s*int64(n-i-1)
			ans[arr[i].idx] = total + left + right
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, ans[i])
		}
		writer.WriteByte('\n')
	}
}

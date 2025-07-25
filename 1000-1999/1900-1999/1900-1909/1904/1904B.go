package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int64
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]Pair, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i].val)
			arr[i].idx = i
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		pref := make([]int64, n)
		pref[0] = arr[0].val
		for i := 1; i < n; i++ {
			pref[i] = pref[i-1] + arr[i].val
		}
		ans := make([]int, n)
		j := 0
		for p := 0; p < n; p++ {
			if j < p {
				j = p
			}
			for j < n-1 && pref[j] >= arr[j+1].val {
				j++
			}
			ans[arr[p].idx] = j
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

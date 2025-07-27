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
		fmt.Fscan(reader, &n)
		u := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &u[i])
		}
		s := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &s[i])
		}

		groups := make(map[int][]int64)
		for i := 0; i < n; i++ {
			groups[u[i]] = append(groups[u[i]], s[i])
		}

		prefixSum := make(map[int][]int64)
		for _, arr := range groups {
			sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
			for i := 1; i < len(arr); i++ {
				arr[i] += arr[i-1]
			}
			size := len(arr)
			ps, ok := prefixSum[size]
			if !ok {
				ps = make([]int64, size)
				prefixSum[size] = ps
			}
			for i := 0; i < size; i++ {
				ps[i] += arr[i]
			}
		}

		ans := make([]int64, n+1)
		for size, ps := range prefixSum {
			for k := 1; k <= size; k++ {
				teams := size / k
				if teams == 0 {
					break
				}
				idx := teams*k - 1
				ans[k] += ps[idx]
			}
		}

		for k := 1; k <= n; k++ {
			if k > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[k])
		}
		fmt.Fprintln(writer)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		A := make([]int, n)
		pos := make([]int, 2*n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &A[i])
			pos[A[i]] = i
		}
		B := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &B[i])
			pos[B[i]] = i
		}

		type pair struct{ val, idx int }
		var arr []pair
		for v := 1; v <= 2*n; v++ {
			arr = append(arr, pair{val: v, idx: pos[v]})
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].val < arr[j].val
		})

		maxIdx := 2 * n
		ok := true
		for i := 2*n - 1; i >= 0; i-- {
			if arr[i].idx > maxIdx {
				ok = false
				break
			}
			if arr[i].idx < maxIdx {
				maxIdx = arr[i].idx
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

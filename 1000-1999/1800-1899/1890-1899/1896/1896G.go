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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var mode string
		fmt.Fscan(in, &mode)
		if mode != "manual" {
			return
		}
		sz := n * n
		perm := make([]int, sz)
		for i := 0; i < sz; i++ {
			fmt.Fscan(in, &perm[i])
		}
		type pair struct{ idx, val int }
		arr := make([]pair, sz)
		for i := 0; i < sz; i++ {
			arr[i] = pair{i + 1, perm[i]}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val > arr[j].val })
		k := sz - n + 1
		fmt.Fprint(out, "!")
		for i := 0; i < k; i++ {
			fmt.Fprintf(out, " %d", arr[i].idx)
		}
		if t > 0 {
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out)
		}
	}
}

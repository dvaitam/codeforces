package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		arr := make([]Pair, n)
		for i := 0; i < n; i++ {
			arr[i] = Pair{a[i], i}
		}
		sort.Slice(arr, func(i, j int) bool {
			if arr[i].val == arr[j].val {
				return arr[i].idx < arr[j].idx
			}
			return arr[i].val > arr[j].val
		})

		processed := make([]bool, n)
		clones := 0
		for _, p := range arr {
			i := p.idx
			left := i > 0 && processed[i-1]
			right := i+1 < n && processed[i+1]
			if !left && !right {
				clones++
			}
			processed[i] = true
		}
		fmt.Fprintln(out, clones)
	}
}

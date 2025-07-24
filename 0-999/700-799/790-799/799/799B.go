package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	prices := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &prices[i])
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	type pair struct{ price, idx int }
	lists := make([][]pair, 4)
	for i := 0; i < n; i++ {
		if a[i] == b[i] {
			lists[a[i]] = append(lists[a[i]], pair{prices[i], i})
		} else {
			lists[a[i]] = append(lists[a[i]], pair{prices[i], i})
			lists[b[i]] = append(lists[b[i]], pair{prices[i], i})
		}
	}
	for c := 1; c <= 3; c++ {
		sort.Slice(lists[c], func(i, j int) bool {
			return lists[c][i].price < lists[c][j].price
		})
	}
	ptr := make([]int, 4)
	sold := make([]bool, n)

	var m int
	fmt.Fscan(in, &m)
	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < m; i++ {
		var c int
		fmt.Fscan(in, &c)
		// advance pointer to next unsold shirt
		for ptr[c] < len(lists[c]) && sold[lists[c][ptr[c]].idx] {
			ptr[c]++
		}
		if ptr[c] == len(lists[c]) {
			fmt.Fprint(out, -1)
		} else {
			pr := lists[c][ptr[c]].price
			sold[lists[c][ptr[c]].idx] = true
			fmt.Fprint(out, pr)
			ptr[c]++
		}
		if i+1 == m {
			fmt.Fprintln(out)
		} else {
			fmt.Fprint(out, " ")
		}
	}
	out.Flush()
}

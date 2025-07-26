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

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	m := 1 << uint(n)

	rank := make([]int, m)
	for i := 0; i < m; i++ {
		rank[i] = int(s[i])
	}
	type pair struct {
		a, b int
		idx  int
	}
	for step := 1; step < m; step <<= 1 {
		arr := make([]pair, m)
		for i := 0; i < m; i++ {
			arr[i] = pair{rank[i], rank[i^step], i}
		}
		sort.Slice(arr, func(i, j int) bool {
			if arr[i].a == arr[j].a {
				return arr[i].b < arr[j].b
			}
			return arr[i].a < arr[j].a
		})
		newRank := make([]int, m)
		cur := 0
		newRank[arr[0].idx] = cur
		for i := 1; i < m; i++ {
			if arr[i].a != arr[i-1].a || arr[i].b != arr[i-1].b {
				cur++
			}
			newRank[arr[i].idx] = cur
		}
		rank = newRank
	}
	best := 0
	for i := 1; i < m; i++ {
		if rank[i] < rank[best] {
			best = i
		}
	}
	res := make([]byte, m)
	for i := 0; i < m; i++ {
		res[i] = s[i^best]
	}
	fmt.Fprintln(out, string(res))
}

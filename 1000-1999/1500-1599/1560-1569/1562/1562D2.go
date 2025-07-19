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
		var n, q int
		fmt.Fscan(reader, &n, &q)
		var s string
		fmt.Fscan(reader, &s)
		p := make([]int, n+1)
		v1 := make([][]int, n+1)
		v2 := make([][]int, n+1)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				if s[i] == '+' {
					p[i+1] = p[i] + 1
				} else {
					p[i+1] = p[i] - 1
				}
			} else {
				if s[i] == '+' {
					p[i+1] = p[i] - 1
				} else {
					p[i+1] = p[i] + 1
				}
			}
			if p[i+1] >= 0 {
				v1[p[i+1]] = append(v1[p[i+1]], i+1)
			} else {
				v2[-p[i+1]] = append(v2[-p[i+1]], i+1)
			}
		}
		for j := 0; j < q; j++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			c := p[r] - p[l-1]
			if c == 0 {
				fmt.Fprintln(writer, 0)
				continue
			}
			if c%2 == 0 {
				// two operations
				// first at original r
				r0 := r
				// compute on shortened range
				r = r0 - 1
				// compute target f
				var f int
				if p[r] > p[l-1] {
					f = p[l-1] + (p[r]-p[l-1]+1)/2
				} else {
					f = p[l-1] - (p[l-1]-p[r]+1)/2
				}
				// find position
				var pos int
				if f >= 0 {
					arr := v1[f]
					idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
					pos = arr[idx]
				} else {
					arr := v2[-f]
					idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
					pos = arr[idx]
				}
				fmt.Fprintln(writer, 2)
				fmt.Fprintf(writer, "%d %d\n", r0, pos)
			} else {
				// one operation
				// compute target f
				var f int
				if p[r] > p[l-1] {
					f = p[l-1] + (p[r]-p[l-1]+1)/2
				} else {
					f = p[l-1] - (p[l-1]-p[r]+1)/2
				}
				var pos int
				if f >= 0 {
					arr := v1[f]
					idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
					pos = arr[idx]
				} else {
					arr := v2[-f]
					idx := sort.Search(len(arr), func(i int) bool { return arr[i] > l-1 })
					pos = arr[idx]
				}
				fmt.Fprintln(writer, 1)
				fmt.Fprintln(writer, pos)
			}
		}
	}
}

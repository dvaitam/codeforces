package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segment struct {
	start, end int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// build segments of equal values
		segs := make([]segment, 0, n)
		for i := 0; i < n; i++ {
			j := i
			for j < n && a[j] == a[i] {
				j++
			}
			// segment covers [i+1, j]
			segs = append(segs, segment{i + 1, j})
			i = j - 1
		}
		var q int
		fmt.Fscan(reader, &q)
		for q > 0 {
			q--
			var l, r int
			fmt.Fscan(reader, &l, &r)
			// find first segment with start > l
			idx := sort.Search(len(segs), func(i int) bool {
				return segs[i].start > l
			})
			idx--
			if idx >= 0 && segs[idx].end < r {
				fmt.Fprintln(writer, l, segs[idx].end+1)
			} else {
				fmt.Fprintln(writer, -1, -1)
			}
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l, r, c int
}

// difficulty computes the minimal number of changes needed
// to make the sequence nice as described in the problem statement.
// It uses a straightforward approach by constructing intervals
// for each distinct value and then merging overlapping ones.
func difficulty(a []int) int {
	first := make(map[int]int)
	last := make(map[int]int)
	cnt := make(map[int]int)
	for i, v := range a {
		if _, ok := first[v]; !ok {
			first[v] = i
		}
		last[v] = i
		cnt[v]++
	}
	intervals := make([]interval, 0, len(first))
	for v := range first {
		intervals = append(intervals, interval{first[v], last[v], cnt[v]})
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].l == intervals[j].l {
			return intervals[i].r < intervals[j].r
		}
		return intervals[i].l < intervals[j].l
	})
	res := 0
	i := 0
	for i < len(intervals) {
		curL := intervals[i].l
		curR := intervals[i].r
		maxC := intervals[i].c
		j := i + 1
		for j < len(intervals) && intervals[j].l <= curR {
			if intervals[j].r > curR {
				curR = intervals[j].r
			}
			if intervals[j].c > maxC {
				maxC = intervals[j].c
			}
			j++
		}
		res += (curR - curL + 1) - maxC
		i = j
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	fmt.Fprintln(writer, difficulty(a))
	for ; q > 0; q-- {
		var idx, x int
		fmt.Fscan(reader, &idx, &x)
		idx--
		a[idx] = x
		fmt.Fprintln(writer, difficulty(a))
	}
}

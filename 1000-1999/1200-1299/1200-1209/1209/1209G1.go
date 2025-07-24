package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for CF problem 1209G1.
// Given an array with q=0, compute minimal number of elements to change
// so that each distinct value forms one contiguous block. If we change
// any value x, all of its occurrences must be changed to the same new value.
// The optimal strategy is to merge all values whose index intervals
// overlap into groups. Every group forms a contiguous segment of the
// array. Inside each group we can make all elements equal to the most
// frequent value of that group. The number of changes for a group is the
// segment length minus that maximal frequency.

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	const maxV = 200000 + 5
	l := make([]int, maxV)
	r := make([]int, maxV)
	cnt := make([]int, maxV)
	for i := range l {
		l[i] = n
		r[i] = -1
	}
	for i, v := range a {
		if l[v] == n {
			l[v] = i
		}
		r[v] = i
		cnt[v]++
	}

	type interval struct {
		l, r int
		c    int
	}
	var intervals []interval
	for v := 1; v < maxV; v++ {
		if cnt[v] > 0 {
			intervals = append(intervals, interval{l[v], r[v], cnt[v]})
		}
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].l == intervals[j].l {
			return intervals[i].r < intervals[j].r
		}
		return intervals[i].l < intervals[j].l
	})

	ans := 0
	i := 0
	for i < len(intervals) {
		curL := intervals[i].l
		curR := intervals[i].r
		maxCnt := intervals[i].c
		j := i + 1
		for j < len(intervals) && intervals[j].l <= curR {
			if intervals[j].r > curR {
				curR = intervals[j].r
			}
			if intervals[j].c > maxCnt {
				maxCnt = intervals[j].c
			}
			j++
		}
		segLen := curR - curL + 1
		ans += segLen - maxCnt
		i = j
	}

	fmt.Println(ans)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	x     int64
	delta int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var T int64
	if _, err := fmt.Fscan(in, &n, &T); err != nil {
		return
	}
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &t[i])
	}
	events := make([]event, 0, 2*n)
	for i := 1; i <= n; i++ {
		ti := t[i-1]
		if ti >= T { // bowl not ready before show ends
			continue
		}
		l := ti - int64(i)
		if l < 0 {
			l = 0
		}
		r := T - 1 - int64(i)
		if r < l {
			continue
		}
		events = append(events, event{l, 1})
		events = append(events, event{r + 1, -1})
	}
	sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })
	cur := 0
	best := 0
	i := 0
	for i < len(events) {
		x := events[i].x
		for i < len(events) && events[i].x == x {
			cur += events[i].delta
			i++
		}
		if cur > best {
			best = cur
		}
	}
	fmt.Println(best)
}

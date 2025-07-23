package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	x   int
	typ int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	events := make([]event, 0, 2*n)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		events = append(events, event{l, 1}, event{r, -1})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].x == events[j].x {
			return events[i].typ > events[j].typ
		}
		return events[i].x < events[j].x
	})

	segments := make([][2]int, 0)
	cover := 0
	start := 0
	inSeg := false
	for _, e := range events {
		if e.typ == 1 {
			cover++
			if cover >= k && !inSeg {
				start = e.x
				inSeg = true
			}
		} else {
			if cover >= k && cover-1 < k {
				segments = append(segments, [2]int{start, e.x})
				inSeg = false
			}
			cover--
		}
	}

	fmt.Fprintln(writer, len(segments))
	for _, seg := range segments {
		fmt.Fprintf(writer, "%d %d\n", seg[0], seg[1])
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	pos   int64
	delta int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		events := make([]event, 0, 2*n)
		for i := 0; i < n; i++ {
			var a, b int64
			fmt.Fscan(reader, &a, &b)
			events = append(events, event{pos: a, delta: 1})
			events = append(events, event{pos: b + 1, delta: -1})
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].pos == events[j].pos {
				return events[i].delta > events[j].delta
			}
			return events[i].pos < events[j].pos
		})
		var count int64
		var prev int64
		var ans int64 = -1
		if len(events) > 0 {
			prev = events[0].pos
		}
		for _, e := range events {
			if count == 1 && e.pos > prev {
				ans = prev
				break
			}
			count += int64(e.delta)
			prev = e.pos
		}
		if ans == -1 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}

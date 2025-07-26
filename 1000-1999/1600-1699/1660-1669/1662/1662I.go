package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	pos   int64
	delta int64
}

// This program solves problemI.txt for contest 1662I.
// It finds the best position to open a new ice cream shop on a line of huts
// with existing shops so that the new shop is strictly closer for as many
// huts as possible. For each hut we compute the interval of positions where
// the new shop would be closest and then determine the point covered by the
// maximum total population using a sweep line over weighted intervals.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	p := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	x := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &x[i])
	}
	sort.Slice(x, func(i, j int) bool { return x[i] < x[j] })

	events := make([]event, 0, 2*n)
	j := 0
	const INF int64 = 1<<62 - 1
	for i := 0; i < n; i++ {
		pos := int64(i) * 100
		for j < m && x[j] <= pos {
			j++
		}
		left := INF
		if j > 0 {
			left = pos - x[j-1]
		}
		right := INF
		if j < m {
			right = x[j] - pos
		}
		d := left
		if right < d {
			d = right
		}
		if d > 0 {
			events = append(events, event{pos - d, p[i]})
			events = append(events, event{pos + d, -p[i]})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].pos == events[j].pos {
			return events[i].delta < events[j].delta
		}
		return events[i].pos < events[j].pos
	})

	var cur, best int64
	for i := 0; i < len(events); {
		pos := events[i].pos
		for i < len(events) && events[i].pos == pos {
			cur += events[i].delta
			i++
		}
		if cur > best {
			best = cur
		}
	}

	fmt.Fprintln(writer, best)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Update struct {
	t int64
	h int64
}

type Enemy struct {
	maxH  int64
	start int64
	regen int64
	ups   []Update
}

type Event struct {
	t     int64
	delta int64
}

func addInterval(l, r int64, events *[]Event) {
	if l > r {
		return
	}
	*events = append(*events, Event{t: l, delta: 1}, Event{t: r + 1, delta: -1})
}

func processSegment(start, end, h, regen, maxH, damage int64, events *[]Event) {
	if start >= end || h > damage {
		return
	}
	if regen == 0 || maxH <= damage {
		addInterval(start, end-1, events)
		return
	}
	if damage < h {
		return
	}
	limit := (damage - h) / regen
	tEnd := start + limit
	if tEnd >= end {
		tEnd = end - 1
	}
	if tEnd >= start {
		addInterval(start, tEnd, events)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var bounty, increase, damage int64
	fmt.Fscan(reader, &bounty, &increase, &damage)

	enemies := make([]Enemy, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &enemies[i].maxH, &enemies[i].start, &enemies[i].regen)
	}
	for i := 0; i < m; i++ {
		var t, id, h int64
		fmt.Fscan(reader, &t, &id, &h)
		id--
		enemies[id].ups = append(enemies[id].ups, Update{t: t, h: h})
	}

	const INF int64 = 1 << 60
	var events []Event
	infinite := false

	for i := 0; i < n; i++ {
		e := &enemies[i]
		sort.Slice(e.ups, func(a, b int) bool { return e.ups[a].t < e.ups[b].t })
		prevT := int64(0)
		prevH := e.start
		for _, u := range e.ups {
			processSegment(prevT, u.t, prevH, e.regen, e.maxH, damage, &events)
			prevT = u.t
			prevH = u.h
		}
		processSegment(prevT, INF, prevH, e.regen, e.maxH, damage, &events)
		if increase > 0 {
			finalH := prevH
			if e.regen > 0 {
				finalH = e.maxH
			}
			if finalH <= damage {
				infinite = true
			}
		}
	}

	if infinite {
		fmt.Fprintln(writer, -1)
		return
	}

	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
	var curr int64
	var ans int64
	prev := int64(0)
	for i := 0; i < len(events); {
		t := events[i].t
		if curr > 0 && prev <= t-1 {
			gold := curr * (bounty + increase*(t-1))
			if gold > ans {
				ans = gold
			}
		}
		for i < len(events) && events[i].t == t {
			curr += events[i].delta
			i++
		}
		prev = t
	}

	fmt.Fprintln(writer, ans)
}

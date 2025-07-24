package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Event represents start or end of availability for Vasya or Petya.
type Event struct {
	pos int64
	typ int // +1/-1 for Vasya start/end, +2/-2 for Petya start/end
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	var C int64
	if _, err := fmt.Fscan(in, &n, &m, &C); err != nil {
		return
	}

	events := make([]Event, 0, 2*(n+m))
	for i := 0; i < n; i++ {
		var a, b int64
		fmt.Fscan(in, &a, &b)
		events = append(events, Event{a, 1}, Event{b, -1})
	}
	for i := 0; i < m; i++ {
		var c, d int64
		fmt.Fscan(in, &c, &d)
		events = append(events, Event{c, 2}, Event{d, -2})
	}

	sort.Slice(events, func(i, j int) bool { return events[i].pos < events[j].pos })

	var vasya, petya bool
	var prev int64
	var diff int64
	var xp int64

	for _, e := range events {
		if e.pos > prev {
			length := e.pos - prev
			switch {
			case vasya && petya:
				// both available
				if diff > C {
					wait := diff - C
					if wait >= length {
						xp += length // play without boost
						// diff unchanged
					} else {
						xp += 2 * (length - wait)
						diff = C
					}
				} else if diff < -C {
					wait := -C - diff
					if wait >= length {
						xp += length
						diff += length
					} else {
						xp += wait
						xp += 2 * (length - wait)
						diff = -C
					}
				} else {
					xp += 2 * length
				}
			case vasya:
				xp += length
				diff += length
			case petya:
				diff -= length
			}
		}
		prev = e.pos
		switch e.typ {
		case 1:
			vasya = true
		case -1:
			vasya = false
		case 2:
			petya = true
		case -2:
			petya = false
		}
	}

	fmt.Fprintln(out, xp)
}

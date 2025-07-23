package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// interval represents a half-open interval [l, r).
type interval struct {
	l, r int
}

// addInterval adds interval iv to a sorted slice of non-overlapping intervals
// and returns the resulting slice of non-overlapping intervals.
func addInterval(ints []interval, iv interval) []interval {
	n := len(ints)
	pos := sort.Search(n, func(i int) bool { return ints[i].l >= iv.l })
	ints = append(ints, interval{})
	copy(ints[pos+1:], ints[pos:])
	ints[pos] = iv

	res := make([]interval, 0, len(ints))
	for _, it := range ints {
		if len(res) == 0 || res[len(res)-1].r < it.l {
			res = append(res, it)
		} else if res[len(res)-1].r < it.r {
			res[len(res)-1].r = it.r
		}
	}
	return res
}

// length calculates the total length covered by intervals.
func length(ints []interval) int {
	sum := 0
	for _, it := range ints {
		sum += it.r - it.l
	}
	return sum
}

func copyIntervals(ints []interval) []interval {
	res := make([]interval, len(ints))
	copy(res, ints)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	type light struct{ a, l int }
	lights := make([]light, n)
	for i := range lights {
		fmt.Fscan(in, &lights[i].a, &lights[i].l)
	}
	sort.Slice(lights, func(i, j int) bool { return lights[i].a < lights[j].a })

	intervals := []interval{}
	for _, lt := range lights {
		north := interval{lt.a, lt.a + lt.l}
		south := interval{lt.a - lt.l, lt.a}

		intsNorth := addInterval(copyIntervals(intervals), north)
		lenNorth := length(intsNorth)
		intsSouth := addInterval(copyIntervals(intervals), south)
		lenSouth := length(intsSouth)

		if lenNorth >= lenSouth {
			intervals = intsNorth
		} else {
			intervals = intsSouth
		}
	}

	fmt.Println(length(intervals))
}

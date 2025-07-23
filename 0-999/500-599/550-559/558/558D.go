package main

import (
	"bufio"
	"fmt"
	"os"
)

type Interval struct {
	l, r int64
}

func intersect(intervals []Interval, seg Interval) []Interval {
	res := make([]Interval, 0, len(intervals))
	for _, in := range intervals {
		l := in.l
		if seg.l > l {
			l = seg.l
		}
		r := in.r
		if seg.r < r {
			r = seg.r
		}
		if l <= r {
			res = append(res, Interval{l, r})
		}
	}
	return res
}

func subtract(intervals []Interval, seg Interval) []Interval {
	res := make([]Interval, 0, len(intervals)+1)
	for _, in := range intervals {
		if seg.r < in.l || seg.l > in.r {
			res = append(res, in)
			continue
		}
		if seg.l > in.l {
			res = append(res, Interval{in.l, seg.l - 1})
		}
		if seg.r < in.r {
			res = append(res, Interval{seg.r + 1, in.r})
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var h, q int
	if _, err := fmt.Fscan(reader, &h, &q); err != nil {
		return
	}
	allowed := []Interval{{1 << (h - 1), (1 << h) - 1}}
	for ; q > 0; q-- {
		var i int
		var L, R int64
		var ans int
		fmt.Fscan(reader, &i, &L, &R, &ans)
		shift := uint(h - i)
		seg := Interval{L << shift, ((R + 1) << shift) - 1}
		if ans == 1 {
			allowed = intersect(allowed, seg)
		} else {
			allowed = subtract(allowed, seg)
		}
	}
	if len(allowed) == 0 {
		fmt.Println("Game cheated!")
		return
	}
	if len(allowed) == 1 && allowed[0].l == allowed[0].r {
		fmt.Println(allowed[0].l)
		return
	}
	// check if union of intervals results in exactly one number
	count := 0
	var val int64
	for _, in := range allowed {
		if in.l == in.r {
			count++
			val = in.l
		} else {
			// there is a range with more than one value
			fmt.Println("Data not sufficient!")
			return
		}
	}
	if count == 1 {
		fmt.Println(val)
	} else {
		fmt.Println("Data not sufficient!")
	}
}

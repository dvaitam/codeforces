package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	var text string
	if _, err := fmt.Fscan(in, &text); err != nil {
		return
	}
	s := []byte(text)
	n := len(s)

	// collect indices where we are allowed to break the line
	breaks := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if s[i] == '-' || s[i] == ' ' {
			breaks = append(breaks, i)
		}
	}
	breaks = append(breaks, n-1) // end of text

	// helper to check if width w is feasible
	can := func(w int) bool {
		lines := 1
		pos := 0
		idx := sort.SearchInts(breaks, pos) // first break >= pos
		for pos < n {
			if n-pos <= w {
				return lines <= k
			}
			limit := pos + w - 1
			for idx < len(breaks) && breaks[idx] <= limit {
				idx++
			}
			j := idx - 1
			if j < 0 || breaks[j] < pos {
				return false
			}
			pos = breaks[j] + 1
			lines++
			if lines > k {
				return false
			}
		}
		return lines <= k
	}

	// lower bound for width is length of longest unbreakable segment
	maxSeg := 0
	last := -1
	for _, b := range breaks {
		seg := b - last
		if seg > maxSeg {
			maxSeg = seg
		}
		last = b
	}

	l, r := maxSeg, n
	for l < r {
		m := (l + r) / 2
		if can(m) {
			r = m
		} else {
			l = m + 1
		}
	}
	fmt.Fprintln(out, l)
}

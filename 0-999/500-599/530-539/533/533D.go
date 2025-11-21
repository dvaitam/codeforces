package main

import (
	"bufio"
	"fmt"
	"os"
)

// segment describes a boundary interval [start, end) where idx is
// the closest surviving column to the boundary from the corresponding side.
type segment struct {
	start int64
	end   int64
	idx   int
}

// hasSurvivor checks whether the original configuration leaves any ordinary
// column standing (no extra column added). It uses the fact that the order of
// crashes does not affect the final set of survivors.
func hasSurvivor(x []int64, d []int64) bool {
	n := len(x) - 2
	prev := make([]int, n+2)
	next := make([]int, n+2)
	for i := 0; i < n+2; i++ {
		prev[i] = i - 1
		next[i] = i + 1
	}
	removed := make([]bool, n+2)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if x[i+1]-x[i-1] > 2*d[i] {
			q = append(q, i)
		}
	}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if removed[v] {
			continue
		}
		removed[v] = true
		l, r := prev[v], next[v]
		next[l] = r
		prev[r] = l
		for _, nb := range []int{l, r} {
			if nb >= 1 && nb <= n && !removed[nb] {
				if x[next[nb]]-x[prev[nb]] > 2*d[nb] {
					q = append(q, nb)
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !removed[i] {
			return true
		}
	}
	return false
}

// buildSegments sweeps the boundary from the left bearing (0) to the right
// bearing, tracking which column is the closest surviving one to the boundary
// on this side. Segments cover [0, total].
func buildSegments(x []int64, d []int64) []segment {
	const inf int64 = 4_000_000_000_000_000_000
	n := len(x) - 2
	total := x[n+1]

	stack := []int{0} // survivors from the left, last element is closest to boundary
	cur := int64(0)
	nextIdx := 1
	segments := make([]segment, 0)

	// restore stability for current boundary position t
	var stabilize func(t int64)
	stabilize = func(t int64) {
		for {
			changed := false
			// the closest column to the boundary must handle the boundary distance
			for len(stack) >= 2 {
				last := stack[len(stack)-1]
				left := stack[len(stack)-2]
				limit := x[left] + 2*d[last]
				if t > limit {
					stack = stack[:len(stack)-1]
					changed = true
				} else {
					break
				}
			}
			if len(stack) < 3 {
				if !changed {
					break
				}
				continue
			}
			mid := stack[len(stack)-2]
			left := stack[len(stack)-3]
			right := stack[len(stack)-1]
			if x[right]-x[left] > 2*d[mid] {
				// mid cannot sustain the gap between its neighbors
				stack = append(stack[:len(stack)-2], stack[len(stack)-1])
				changed = true
				continue
			}
			if !changed {
				break
			}
		}
	}

	stabilize(cur)

	for cur < total {
		var limit int64
		if len(stack) == 1 {
			limit = inf
		} else {
			last := stack[len(stack)-1]
			left := stack[len(stack)-2]
			limit = x[left] + 2*d[last]
		}

		nextAdd := total
		if nextIdx <= n {
			nextAdd = x[nextIdx]
		}

		event := limit
		if nextAdd < event {
			event = nextAdd
		}

		segments = append(segments, segment{start: cur, end: event, idx: stack[len(stack)-1]})
		cur = event

		if cur == nextAdd && nextIdx <= n {
			left := stack[len(stack)-1]
			if cur-x[left] <= 2*d[nextIdx] {
				stack = append(stack, nextIdx)
			}
			stabilize(cur)
			nextIdx++
		} else if cur == limit && limit < nextAdd {
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			stabilize(cur)
		} else if cur == nextAdd && cur == limit && nextIdx <= n {
			left := stack[len(stack)-1]
			if cur-x[left] <= 2*d[nextIdx] {
				stack = append(stack, nextIdx)
			}
			stabilize(cur)
			nextIdx++
		}
	}

	// merge adjacent intervals with the same idx to simplify later merging
	merged := make([]segment, 0, len(segments))
	for _, s := range segments {
		if len(merged) > 0 && merged[len(merged)-1].idx == s.idx && merged[len(merged)-1].end == s.start {
			merged[len(merged)-1].end = s.end
		} else {
			merged = append(merged, s)
		}
	}
	return merged
}

// minRequired returns the minimal width between nearest left and right survivors
// over all boundary positions. The needed durability is this width divided by two.
func minRequired(x []int64, d []int64) float64 {
	leftSeg := buildSegments(x, d)

	n := len(x) - 2
	total := x[n+1]
	xr := make([]int64, n+2)
	dr := make([]int64, n+2)
	xr[0] = 0
	xr[n+1] = total
	for i := 1; i <= n; i++ {
		xr[i] = total - x[n+1-i]
		dr[i] = d[n+1-i]
	}

	rightMirror := buildSegments(xr, dr)
	rightSeg := make([]segment, 0, len(rightMirror))
	for i := len(rightMirror) - 1; i >= 0; i-- {
		s := rightMirror[i]
		start := total - s.end
		end := total - s.start
		idxMir := s.idx
		idxOrig := 0
		if idxMir == 0 {
			idxOrig = n + 1
		} else {
			idxOrig = n + 1 - idxMir
		}
		rightSeg = append(rightSeg, segment{start: start, end: end, idx: idxOrig})
	}

	mergedR := make([]segment, 0, len(rightSeg))
	for _, s := range rightSeg {
		if len(mergedR) > 0 && mergedR[len(mergedR)-1].idx == s.idx && mergedR[len(mergedR)-1].end == s.start {
			mergedR[len(mergedR)-1].end = s.end
		} else {
			mergedR = append(mergedR, s)
		}
	}
	rightSeg = mergedR

	ans := float64(total)
	i, j := 0, 0
	for i < len(leftSeg) && j < len(rightSeg) {
		ls := leftSeg[i]
		rs := rightSeg[j]
		start := ls.start
		if rs.start > start {
			start = rs.start
		}
		end := ls.end
		if rs.end < end {
			end = rs.end
		}
		if start < end {
			width := float64(x[rs.idx] - x[ls.idx])
			if width < ans {
				ans = width
			}
		}
		if ls.end < rs.end {
			i++
		} else if ls.end > rs.end {
			j++
		} else {
			i++
			j++
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	x := make([]int64, n+2)
	for i := 0; i < n+2; i++ {
		fmt.Fscan(in, &x[i])
	}
	d := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &d[i])
	}

	if hasSurvivor(x, d) {
		fmt.Printf("0\n")
		return
	}

	width := minRequired(x, d)
	fmt.Printf("%.10f\n", width/2.0)
}

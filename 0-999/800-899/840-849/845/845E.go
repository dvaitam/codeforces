package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	x   int64
	idx int
	typ int
}

func can(n, m int64, x, y []int64, r int64) bool {
	k := len(x)
	events := make([]event, 0, 2*k+1)
	for i := 0; i < k; i++ {
		lx := x[i] - r
		rx := x[i] + r
		if rx < 1 || lx > n {
			continue
		}
		if lx < 1 {
			lx = 1
		}
		if rx > n {
			rx = n
		}
		events = append(events, event{lx, i, 1})
		events = append(events, event{rx + 1, i, -1})
	}
	events = append(events, event{n + 1, -1, 0})
	sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })

	active := make([]bool, k)
	minX, maxX := n+1, int64(0)
	minY, maxY := m+1, int64(0)
	ptr := 0
	xPrev := int64(1)
	for ptr < len(events) {
		xPos := events[ptr].x
		if xPrev < xPos {
			// process segment [xPrev, xPos-1]
			intervals := make([][2]int64, 0, k)
			for i := 0; i < k; i++ {
				if !active[i] {
					continue
				}
				l := y[i] - r
				r2 := y[i] + r
				if r2 < 1 || l > m {
					continue
				}
				if l < 1 {
					l = 1
				}
				if r2 > m {
					r2 = m
				}
				intervals = append(intervals, [2]int64{l, r2})
			}
			if len(intervals) == 0 {
				// entire column uncovered
				if minX > xPrev {
					minX = xPrev
				}
				if maxX < xPos-1 {
					maxX = xPos - 1
				}
				if minY > 1 {
					minY = 1
				}
				if maxY < m {
					maxY = m
				}
			} else {
				sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
				curr := int64(1)
				for _, in := range intervals {
					l, r2 := in[0], in[1]
					if r2 < curr {
						continue
					}
					if l > curr {
						if minX > xPrev {
							minX = xPrev
						}
						if maxX < xPos-1 {
							maxX = xPos - 1
						}
						if minY > curr {
							minY = curr
						}
						if maxY < l-1 {
							maxY = l - 1
						}
					}
					if l <= curr {
						if r2+1 > curr {
							curr = r2 + 1
						}
					} else {
						curr = r2 + 1
					}
					if curr > m {
						break
					}
				}
				if curr <= m {
					if minX > xPrev {
						minX = xPrev
					}
					if maxX < xPos-1 {
						maxX = xPos - 1
					}
					if minY > curr {
						minY = curr
					}
					if maxY < m {
						maxY = m
					}
				}
			}
		}
		for ptr < len(events) && events[ptr].x == xPos {
			if events[ptr].idx >= 0 {
				if events[ptr].typ == 1 {
					active[events[ptr].idx] = true
				} else {
					active[events[ptr].idx] = false
				}
			}
			ptr++
		}
		xPrev = xPos
	}

	if maxX == 0 {
		return true
	}
	width := maxX - minX + 1
	height := maxY - minY + 1
	if width <= 2*r+1 && height <= 2*r+1 {
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int64
	var k int
	fmt.Fscan(reader, &n, &m, &k)
	x := make([]int64, k)
	y := make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}
	lo, hi := int64(0), int64(0)
	if n > m {
		hi = n
	} else {
		hi = m
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if can(n, m, x, y, mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	fmt.Println(lo)
}

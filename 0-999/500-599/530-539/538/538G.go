package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, l int
	if _, err := fmt.Fscan(in, &n, &l); err != nil {
		return
	}

	times := make([]int64, n)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &times[i], &xs[i], &ys[i])
	}

	aVals := make([]int64, n)
	bVals := make([]int64, n)
	for i := 0; i < n; i++ {
		aVals[i] = xs[i] + ys[i]
		bVals[i] = xs[i] - ys[i]
	}

	aSteps, ok := solve1D(l, times, aVals)
	if !ok {
		fmt.Println("NO")
		return
	}
	bSteps, ok := solve1D(l, times, bVals)
	if !ok {
		fmt.Println("NO")
		return
	}

	var sb strings.Builder
	sb.Grow(l)
	for i := 0; i < l; i++ {
		da := aSteps[i]
		db := bSteps[i]
		var ch byte
		switch {
		case da == 1 && db == 1:
			ch = 'R'
		case da == 1 && db == -1:
			ch = 'U'
		case da == -1 && db == 1:
			ch = 'D'
		case da == -1 && db == -1:
			ch = 'L'
		default:
			fmt.Println("NO")
			return
		}
		sb.WriteByte(ch)
	}
	fmt.Println(sb.String())
}

type point struct {
	idx         int
	alpha, beta int64
}

func solve1D(l int, times []int64, values []int64) ([]int8, bool) {
	l64 := int64(l)
	seen := make([]bool, l)
	firstQ := make([]int64, l)
	firstVal := make([]int64, l)
	seen[0] = true
	firstQ[0] = 0
	firstVal[0] = 0

	var T int64
	tDetermined := false

	for i := 0; i < len(times); i++ {
		t := times[i]
		v := values[i]
		if ((v - t) & 1) != 0 {
			return nil, false
		}
		q := t / l64
		r := int(t % l64)
		if !seen[r] {
			seen[r] = true
			firstQ[r] = q
			firstVal[r] = v
		} else {
			dq := q - firstQ[r]
			dv := v - firstVal[r]
			if dq == 0 {
				if dv != 0 {
					return nil, false
				}
			} else {
				if dv%dq != 0 {
					return nil, false
				}
				cand := dv / dq
				if tDetermined {
					if cand != T {
						return nil, false
					}
				} else {
					T = cand
					tDetermined = true
				}
			}
		}
	}

	parityT := int64(l & 1)
	if tDetermined {
		if (T & 1) != parityT {
			return nil, false
		}
	}

	points := make([]point, 0, len(times)+2)
	for r := 0; r < l; r++ {
		if seen[r] {
			points = append(points, point{idx: r, alpha: firstVal[r], beta: firstQ[r]})
		}
	}
	points = append(points, point{idx: l, alpha: 0, beta: -1})

	low := -int64(l)
	high := int64(l)

	for i := 0; i+1 < len(points); i++ {
		d := int64(points[i+1].idx - points[i].idx)
		A := points[i+1].alpha - points[i].alpha
		B := points[i+1].beta - points[i].beta

		if B%2 == 0 {
			if ((A - d) & 1) != 0 {
				return nil, false
			}
		} else {
			if ((A - d) & 1) != parityT {
				return nil, false
			}
		}

		if B == 0 {
			if abs64(A) > d {
				return nil, false
			}
			continue
		}

		left := A - d
		right := A + d
		var lBound, rBound int64
		if B > 0 {
			lBound = ceilDiv(left, B)
			rBound = floorDiv(right, B)
		} else {
			lBound = ceilDiv(right, B)
			rBound = floorDiv(left, B)
		}
		if lBound > rBound {
			return nil, false
		}
		if lBound > low {
			low = lBound
		}
		if rBound < high {
			high = rBound
		}
		if low > high {
			return nil, false
		}
	}

	if !tDetermined {
		cand := low
		if (cand & 1) != parityT {
			cand++
		}
		if cand > high {
			return nil, false
		}
		T = cand
	} else {
		if T < low || T > high {
			return nil, false
		}
	}

	positions := make([]int64, l+1)
	known := make([]bool, l+1)
	positions[0] = 0
	known[0] = true
	positions[l] = T
	known[l] = true
	for r := 1; r < l; r++ {
		if seen[r] {
			positions[r] = firstVal[r] - firstQ[r]*T
			known[r] = true
		}
	}

	steps := make([]int8, l)
	prevIdx := 0
	prevVal := positions[0]
	for i := 1; i <= l; i++ {
		if known[i] {
			if !fillSegment(steps, prevIdx, i, prevVal, positions[i]) {
				return nil, false
			}
			prevIdx = i
			prevVal = positions[i]
		}
	}

	return steps, true
}

func fillSegment(steps []int8, startIdx, endIdx int, startVal, endVal int64) bool {
	diff := endVal - startVal
	idx := startIdx
	if diff > 0 {
		for cnt := int(diff); cnt > 0; cnt-- {
			steps[idx] = 1
			idx++
		}
	} else if diff < 0 {
		for cnt := int(-diff); cnt > 0; cnt-- {
			steps[idx] = -1
			idx++
		}
	}
	rem := endIdx - idx
	if rem < 0 || rem%2 != 0 {
		return false
	}
	for rem > 0 {
		steps[idx] = 1
		idx++
		steps[idx] = -1
		idx++
		rem -= 2
	}
	return true
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func floorDiv(a, b int64) int64 {
	if b < 0 {
		a = -a
		b = -b
	}
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func ceilDiv(a, b int64) int64 {
	return -floorDiv(-a, b)
}

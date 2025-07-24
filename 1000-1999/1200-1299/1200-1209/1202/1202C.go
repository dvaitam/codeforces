package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// solve computes minimal possible grid area after inserting at most
// one move among 'W','A','S','D' into s.
func solve(s string) int64 {
	n := len(s)
	x := make([]int, n+1)
	y := make([]int, n+1)
	for i, c := range s {
		x[i+1] = x[i]
		y[i+1] = y[i]
		switch c {
		case 'W':
			y[i+1]++
		case 'S':
			y[i+1]--
		case 'A':
			x[i+1]--
		case 'D':
			x[i+1]++
		}
	}
	prefMinX := make([]int, n+1)
	prefMaxX := make([]int, n+1)
	prefMinY := make([]int, n+1)
	prefMaxY := make([]int, n+1)
	prefMinX[0], prefMaxX[0] = x[0], x[0]
	prefMinY[0], prefMaxY[0] = y[0], y[0]
	for i := 1; i <= n; i++ {
		prefMinX[i] = min(prefMinX[i-1], x[i])
		prefMaxX[i] = max(prefMaxX[i-1], x[i])
		prefMinY[i] = min(prefMinY[i-1], y[i])
		prefMaxY[i] = max(prefMaxY[i-1], y[i])
	}
	sufMinX := make([]int, n+1)
	sufMaxX := make([]int, n+1)
	sufMinY := make([]int, n+1)
	sufMaxY := make([]int, n+1)
	sufMinX[n], sufMaxX[n] = x[n], x[n]
	sufMinY[n], sufMaxY[n] = y[n], y[n]
	for i := n - 1; i >= 0; i-- {
		sufMinX[i] = min(x[i], sufMinX[i+1])
		sufMaxX[i] = max(x[i], sufMaxX[i+1])
		sufMinY[i] = min(y[i], sufMinY[i+1])
		sufMaxY[i] = max(y[i], sufMaxY[i+1])
	}

	widthOrig := prefMaxX[n] - prefMinX[n] + 1
	heightOrig := prefMaxY[n] - prefMinY[n] + 1
	ans := int64(widthOrig) * int64(heightOrig)

	// try inserting horizontal moves
	minWidth := widthOrig
	for _, dx := range []int{-1, 1} {
		for i := 0; i <= n; i++ {
			minVal := prefMinX[i]
			maxVal := prefMaxX[i]
			v := x[i] + dx
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
			if i != n {
				v = sufMinX[i+1] + dx
				if v < minVal {
					minVal = v
				}
				v = sufMaxX[i+1] + dx
				if v > maxVal {
					maxVal = v
				}
			}
			width := maxVal - minVal + 1
			if width < minWidth {
				minWidth = width
			}
		}
	}
	area := int64(minWidth) * int64(heightOrig)
	if area < ans {
		ans = area
	}

	// try inserting vertical moves
	minHeight := heightOrig
	for _, dy := range []int{-1, 1} {
		for i := 0; i <= n; i++ {
			minVal := prefMinY[i]
			maxVal := prefMaxY[i]
			v := y[i] + dy
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
			if i != n {
				v = sufMinY[i+1] + dy
				if v < minVal {
					minVal = v
				}
				v = sufMaxY[i+1] + dy
				if v > maxVal {
					maxVal = v
				}
			}
			height := maxVal - minVal + 1
			if height < minHeight {
				minHeight = height
			}
		}
	}
	area2 := int64(minHeight) * int64(widthOrig)
	if area2 < ans {
		ans = area2
	}

	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solve(s))
	}
}

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type Point struct {
	x  int64
	y  int64
	yi int
}

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, delta int) {
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *Fenwick) Kth(k int) int {
	if k <= 0 {
		return 1
	}
	idx := 0
	bitMask := 1 << bits.Len(uint(f.n))
	for bitMask > 0 {
		next := idx + bitMask
		if next <= f.n && f.bit[next] < k {
			k -= f.bit[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

func subsetInterval(tree *Fenwick, size, k int, yVals []int64) (int64, int64, bool) {
	if size < 2*k {
		return 0, 0, false
	}
	lowerIdx := tree.Kth(k)
	upperIdx := tree.Kth(size - k + 1)
	lower := yVals[lowerIdx-1]
	upper := yVals[upperIdx-1]
	if lower >= upper {
		return 0, 0, false
	}
	return lower, upper, true
}

func findY(left *Fenwick, leftSize int, right *Fenwick, rightSize int, k int, yVals []int64) (bool, int64) {
	lowerL, upperL, ok := subsetInterval(left, leftSize, k, yVals)
	if !ok {
		return false, 0
	}
	lowerR, upperR, ok := subsetInterval(right, rightSize, k, yVals)
	if !ok {
		return false, 0
	}
	low := lowerL
	if lowerR > low {
		low = lowerR
	}
	high := upperL
	if upperR < high {
		high = upperR
	}
	if high > low {
		return true, high
	}
	return false, 0
}

func check(points []Point, yVals []int64, k int) (bool, int64, int64) {
	if k == 0 {
		return true, 0, 0
	}
	m := len(yVals)
	left := NewFenwick(m)
	right := NewFenwick(m)
	for _, p := range points {
		right.Add(p.yi, 1)
	}
	leftSize := 0
	rightSize := len(points)

	i := 0
	for i < len(points) {
		xVal := points[i].x
		if leftSize >= 2*k && rightSize >= 2*k {
			if ok, yVal := findY(left, leftSize, right, rightSize, k, yVals); ok {
				return true, xVal, yVal
			}
		}
		j := i
		for j < len(points) && points[j].x == xVal {
			idx := points[j].yi
			left.Add(idx, 1)
			right.Add(idx, -1)
			leftSize++
			rightSize--
			j++
		}
		i = j
	}
	return false, 0, 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		points := make([]Point, n)
		yVals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &points[i].x, &points[i].y)
			yVals[i] = points[i].y
		}
		sort.Slice(yVals, func(i, j int) bool { return yVals[i] < yVals[j] })
		uniqY := yVals[:0]
		for _, v := range yVals {
			if len(uniqY) == 0 || uniqY[len(uniqY)-1] != v {
				uniqY = append(uniqY, v)
			}
		}
		yVals = uniqY

		yIndex := make(map[int64]int, len(yVals))
		for idx, val := range yVals {
			yIndex[val] = idx + 1
		}
		for i := range points {
			points[i].yi = yIndex[points[i].y]
		}

		sort.Slice(points, func(i, j int) bool {
			if points[i].x == points[j].x {
				return points[i].y < points[j].y
			}
			return points[i].x < points[j].x
		})

		low, high := 0, n/4
		bestK := 0
		var bestX, bestY int64
		for low <= high {
			mid := (low + high) / 2
			ok, xVal, yVal := check(points, yVals, mid)
			if ok {
				bestK = mid
				bestX = xVal
				bestY = yVal
				low = mid + 1
			} else {
				high = mid - 1
			}
		}

		fmt.Fprintln(out, bestK)
		fmt.Fprintf(out, "%d %d\n", bestX, bestY)
	}
}

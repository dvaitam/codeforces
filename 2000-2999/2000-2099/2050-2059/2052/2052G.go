package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type Coord struct {
	a int64
	b int64
}

type Point struct {
	x Coord
	y Coord
}

type Segment struct {
	start Point
	end   Point
}

var (
	dirDx = []Coord{
		{2, 0},  // 0°
		{0, 1},  // 45°
		{0, 0},  // 90°
		{0, -1}, // 135°
		{-2, 0}, // 180°
		{0, -1}, // 225°
		{0, 0},  // 270°
		{0, 1},  // 315°
	}
	dirDy = []Coord{
		{0, 0},
		{0, 1},
		{2, 0},
		{0, 1},
		{0, 0},
		{0, -1},
		{-2, 0},
		{0, -1},
	}
	cosTypes = []int{1, 2, 0, -2, -1, -2, 0, 2}
	sinTypes = []int{0, 2, 1, 2, 0, -2, -1, -2}
)

func addCoord(a, b Coord) Coord {
	return Coord{a.a + b.a, a.b + b.b}
}

func subCoord(a, b Coord) Coord {
	return Coord{a.a - b.a, a.b - b.b}
}

func scaleCoord(c Coord, k int64) Coord {
	return Coord{c.a * k, c.b * k}
}

func addPoint(p Point, dx, dy Coord) Point {
	return Point{
		x: addCoord(p.x, dx),
		y: addCoord(p.y, dy),
	}
}

func mulCoord(c Coord, typ int) Coord {
	switch typ {
	case 0:
		return Coord{}
	case 1:
		return c
	case -1:
		return Coord{-c.a, -c.b}
	case 2:
		return Coord{c.b, c.a / 2}
	case -2:
		return Coord{-c.b, -c.a / 2}
	default:
		panic("invalid type")
	}
}

func rotatePoint(p Point, idx int) Point {
	ct := cosTypes[idx]
	st := sinTypes[idx]
	xNew := subCoord(mulCoord(p.x, ct), mulCoord(p.y, st))
	yNew := addCoord(mulCoord(p.x, st), mulCoord(p.y, ct))
	return Point{x: xNew, y: yNew}
}

func cmp128(hi1, lo1, hi2, lo2 uint64) int {
	if hi1 < hi2 {
		return -1
	}
	if hi1 > hi2 {
		return 1
	}
	if lo1 < lo2 {
		return -1
	}
	if lo1 > lo2 {
		return 1
	}
	return 0
}

func mulSquareTimes2(x uint64) (uint64, uint64) {
	hi, lo := bits.Mul64(x, x)
	carry := lo >> 63
	lo <<= 1
	hi = (hi << 1) | carry
	return hi, lo
}

func cmpSqrt(absdx, abss uint64) int {
	hi1, lo1 := bits.Mul64(absdx, absdx)
	hi2, lo2 := mulSquareTimes2(abss)
	return cmp128(hi1, lo1, hi2, lo2)
}

func compareCoord(a, b Coord) int {
	dx := a.a - b.a
	dy := a.b - b.b
	if dy == 0 {
		switch {
		case dx < 0:
			return -1
		case dx > 0:
			return 1
		default:
			return 0
		}
	}
	if dy > 0 {
		if dx >= 0 {
			return 1
		}
		absdx := uint64(-dx)
		abss := uint64(dy)
		if cmpSqrt(absdx, abss) < 0 {
			return 1
		}
		return -1
	}
	// dy < 0
	if dx <= 0 {
		return -1
	}
	absdx := uint64(dx)
	abss := uint64(-dy)
	if cmpSqrt(absdx, abss) > 0 {
		return 1
	}
	return -1
}

func lessPoint(a, b Point) bool {
	if cmp := compareCoord(a.x, b.x); cmp != 0 {
		return cmp < 0
	}
	return compareCoord(a.y, b.y) < 0
}

func subPoint(a, b Point) Point {
	return Point{
		x: subCoord(a.x, b.x),
		y: subCoord(a.y, b.y),
	}
}

func canonicalize(segs []Segment) []Segment {
	if len(segs) == 0 {
		return nil
	}
	arr := make([]Segment, len(segs))
	copy(arr, segs)
	minPt := arr[0].start
	for _, s := range arr {
		if lessPoint(s.start, minPt) {
			minPt = s.start
		}
		if lessPoint(s.end, minPt) {
			minPt = s.end
		}
	}
	for i := range arr {
		arr[i].start = subPoint(arr[i].start, minPt)
		arr[i].end = subPoint(arr[i].end, minPt)
	}
	sort.Slice(arr, func(i, j int) bool {
		a, b := arr[i], arr[j]
		if lessPoint(a.start, b.start) {
			return true
		}
		if lessPoint(b.start, a.start) {
			return false
		}
		return lessPoint(a.end, b.end)
	})
	return arr
}

func segmentsEqual(a, b []Segment) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].start.x != b[i].start.x || a[i].start.y != b[i].start.y {
			return false
		}
		if a[i].end.x != b[i].end.x || a[i].end.y != b[i].end.y {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	cur := Point{}
	dir := 0
	segments := make([]Segment, 0, 2048)

	for i := 0; i < n; i++ {
		var cmd string
		fmt.Fscan(reader, &cmd)
		switch cmd {
		case "rotate":
			var a int
			fmt.Fscan(reader, &a)
			dir = ((dir+(a/45))%8 + 8) % 8
		case "draw":
			var d int64
			fmt.Fscan(reader, &d)
			start := cur
			dx := scaleCoord(dirDx[dir], d)
			dy := scaleCoord(dirDy[dir], d)
			cur = addPoint(cur, dx, dy)
			seg := Segment{start: start, end: cur}
			if lessPoint(seg.end, seg.start) {
				seg.start, seg.end = seg.end, seg.start
			}
			segments = append(segments, seg)
		case "move":
			var d int64
			fmt.Fscan(reader, &d)
			dx := scaleCoord(dirDx[dir], d)
			dy := scaleCoord(dirDy[dir], d)
			cur = addPoint(cur, dx, dy)
		}
	}

	base := canonicalize(segments)
	answer := 360
	for k := 1; k < 8; k++ {
		rotated := make([]Segment, len(segments))
		for i, seg := range segments {
			rotated[i] = Segment{
				start: rotatePoint(seg.start, k),
				end:   rotatePoint(seg.end, k),
			}
			if lessPoint(rotated[i].end, rotated[i].start) {
				rotated[i].start, rotated[i].end = rotated[i].end, rotated[i].start
			}
		}
		canon := canonicalize(rotated)
		if segmentsEqual(canon, base) {
			answer = k * 45
			break
		}
	}
	fmt.Fprintln(writer, answer)
}

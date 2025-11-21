package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	dirRight = 0
	dirLeft  = 1
	dirUp    = 2
	dirDown  = 3
)

var dirDx = []int{1, -1, 0, 0}
var dirDy = []int{0, 0, 1, -1}

type interval struct {
	l   int
	r   int
	dir int
	id  int
}

type segmentTree struct {
	size int
	tree [][]int
}

func newSegmentTree(limit int) *segmentTree {
	size := 1
	for size < limit+1 {
		size <<= 1
	}
	return &segmentTree{
		size: size,
		tree: make([][]int, size<<1),
	}
}

func (st *segmentTree) addRange(l, r, val int) {
	if l > r {
		return
	}
	st.add(1, 0, st.size-1, l, r, val)
}

func (st *segmentTree) add(node, left, right, ql, qr, val int) {
	if ql <= left && right <= qr {
		st.tree[node] = append(st.tree[node], val)
		return
	}
	mid := (left + right) >> 1
	if ql <= mid {
		st.add(node<<1, left, mid, ql, qr, val)
	}
	if qr > mid {
		st.add(node<<1|1, mid+1, right, ql, qr, val)
	}
}

func (st *segmentTree) build() {
	for i := range st.tree {
		if len(st.tree[i]) > 1 {
			sort.Ints(st.tree[i])
		} else if len(st.tree[i]) == 1 {
			// nothing
		}
	}
}

func (st *segmentTree) query(point, current, dir int) (int, bool) {
	if point < 0 || point >= st.size {
		return 0, false
	}
	best := 0
	found := false
	if dir > 0 {
		best = int(1e9)
	} else {
		best = -int(1e9)
	}
	node := 1
	left, right := 0, st.size-1
	for {
		arr := st.tree[node]
		if len(arr) > 0 {
			if dir > 0 {
				pos := sort.Search(len(arr), func(i int) bool { return arr[i] > current })
				if pos < len(arr) {
					val := arr[pos]
					if !found || val < best {
						best = val
						found = true
					}
				}
			} else {
				pos := sort.Search(len(arr), func(i int) bool { return arr[i] >= current })
				pos--
				if pos >= 0 {
					val := arr[pos]
					if !found || val > best {
						best = val
						found = true
					}
				}
			}
		}
		if left == right {
			break
		}
		mid := (left + right) >> 1
		if point <= mid {
			node = node << 1
			right = mid
		} else {
			node = node<<1 | 1
			left = mid + 1
		}
	}
	return best, found
}

type solver struct {
	b             int
	horizLines    [][]interval
	vertLines     [][]interval
	verticalTree  *segmentTree
	horizontalTree *segmentTree
	nextArrowID   int
}

func newSolver(b int) *solver {
	return &solver{
		b:              b,
		horizLines:     make([][]interval, b+1),
		vertLines:      make([][]interval, b+1),
		verticalTree:   newSegmentTree(b),
		horizontalTree: newSegmentTree(b),
	}
}

func (s *solver) addArrow(x0, y0, x1, y1 int) {
	if x0 == x1 {
		x := x0
		l, r := y0, y1
		if l > r {
			l, r = r, l
		}
		dir := dirDown
		if y1 > y0 {
			dir = dirUp
		}
		s.vertLines[x] = append(s.vertLines[x], interval{l: l, r: r, dir: dir, id: s.nextArrowID})
		s.verticalTree.addRange(l, r, x)
	} else {
		y := y0
		l, r := x0, x1
		if l > r {
			l, r = r, l
		}
		dir := dirLeft
		if x1 > x0 {
			dir = dirRight
		}
		s.horizLines[y] = append(s.horizLines[y], interval{l: l, r: r, dir: dir, id: s.nextArrowID})
		s.horizontalTree.addRange(l, r, y)
	}
	s.nextArrowID++
}

func (s *solver) build() {
	for i := range s.horizLines {
		if len(s.horizLines[i]) > 0 {
			sort.Slice(s.horizLines[i], func(a, b int) bool {
				if s.horizLines[i][a].l == s.horizLines[i][b].l {
					return s.horizLines[i][a].r < s.horizLines[i][b].r
				}
				return s.horizLines[i][a].l < s.horizLines[i][b].l
			})
		}
	}
	for i := range s.vertLines {
		if len(s.vertLines[i]) > 0 {
			sort.Slice(s.vertLines[i], func(a, b int) bool {
				if s.vertLines[i][a].l == s.vertLines[i][b].l {
					return s.vertLines[i][a].r < s.vertLines[i][b].r
				}
				return s.vertLines[i][a].l < s.vertLines[i][b].l
			})
		}
	}
	s.verticalTree.build()
	s.horizontalTree.build()
}

func (s *solver) findHorizontalAt(y, x int) *interval {
	if y < 0 || y > s.b {
		return nil
	}
	line := s.horizLines[y]
	if len(line) == 0 {
		return nil
	}
	idx := sort.Search(len(line), func(i int) bool { return line[i].l > x })
	if idx == 0 {
		return nil
	}
	iv := &line[idx-1]
	if x >= iv.l && x <= iv.r {
		return iv
	}
	return nil
}

func (s *solver) findVerticalAt(x, y int) *interval {
	if x < 0 || x > s.b {
		return nil
	}
	line := s.vertLines[x]
	if len(line) == 0 {
		return nil
	}
	idx := sort.Search(len(line), func(i int) bool { return line[i].l > y })
	if idx == 0 {
		return nil
	}
	iv := &line[idx-1]
	if y >= iv.l && y <= iv.r {
		return iv
	}
	return nil
}

func (s *solver) processArrow(x, y *int, dir *int, lastID *int) bool {
	if iv := s.findHorizontalAt(*y, *x); iv != nil {
		if *lastID == iv.id {
			return false
		}
		*dir = iv.dir
		*lastID = iv.id
		return true
	}
	if iv := s.findVerticalAt(*x, *y); iv != nil {
		if *lastID == iv.id {
			return false
		}
		*dir = iv.dir
		*lastID = iv.id
		return true
	}
	if *lastID != -1 {
		*lastID = -1
	}
	return false
}

func (s *solver) nextHorizontalDistance(x, y, dir int) (int64, bool) {
	var best int
	var boundary bool
	if dir == dirRight {
		best = s.b - x
	} else {
		best = x
	}
	boundary = true

	if val, ok := s.verticalTree.query(y, x, dirDx[dir]); ok {
		dist := val - x
		if dir == dirLeft {
			dist = x - val
		}
		if dist < best || (dist == best && boundary) {
			best = dist
			boundary = false
		}
	}

	line := s.horizLines[y]
	if len(line) > 0 {
		if dir == dirRight {
			idx := sort.Search(len(line), func(i int) bool { return line[i].l > x })
			if idx < len(line) {
				dist := line[idx].l - x
				if dist < best || (dist == best && boundary) {
					best = dist
					boundary = false
				}
			}
		} else {
			idx := sort.Search(len(line), func(i int) bool { return line[i].r >= x })
			idx--
			if idx >= 0 {
				dist := x - line[idx].r
				if dist < best || (dist == best && boundary) {
					best = dist
					boundary = false
				}
			}
		}
	}
	return int64(best), boundary
}

func (s *solver) nextVerticalDistance(x, y, dir int) (int64, bool) {
	var best int
	var boundary bool
	if dir == dirUp {
		best = s.b - y
	} else {
		best = y
	}
	boundary = true

	if val, ok := s.horizontalTree.query(x, y, dirDy[dir]); ok {
		dist := val - y
		if dir == dirDown {
			dist = y - val
		}
		if dist < best || (dist == best && boundary) {
			best = dist
			boundary = false
		}
	}

	line := s.vertLines[x]
	if len(line) > 0 {
		if dir == dirUp {
			idx := sort.Search(len(line), func(i int) bool { return line[i].l > y })
			if idx < len(line) {
				dist := line[idx].l - y
				if dist < best || (dist == best && boundary) {
					best = dist
					boundary = false
				}
			}
		} else {
			idx := sort.Search(len(line), func(i int) bool { return line[i].r >= y })
			idx--
			if idx >= 0 {
				dist := y - line[idx].r
				if dist < best || (dist == best && boundary) {
					best = dist
					boundary = false
				}
			}
		}
	}
	return int64(best), boundary
}

func (s *solver) nextDistance(x, y, dir int) (int64, bool) {
	if dir == dirRight || dir == dirLeft {
		return s.nextHorizontalDistance(x, y, dir)
	}
	return s.nextVerticalDistance(x, y, dir)
}

func encodeState(x, y, dir int) uint64 {
	return (uint64(x) << 32) | (uint64(y) << 2) | uint64(dir)
}

func (s *solver) query(x, y, dir int, t int64) (int, int) {
	posX, posY := x, y
	timeLeft := t
	elapsed := int64(0)
	lastArrowID := -1
	seen := make(map[uint64]int64)

	for timeLeft > 0 {
		for {
			if !s.processArrow(&posX, &posY, &dir, &lastArrowID) {
				break
			}
		}

		state := encodeState(posX, posY, dir)
		if prev, ok := seen[state]; ok {
			cycle := elapsed - prev
			if cycle > 0 {
				skip := timeLeft / cycle
				if skip > 0 {
					elapsed += skip * cycle
					timeLeft -= skip * cycle
					continue
				}
			}
		} else {
			seen[state] = elapsed
		}

		dist, boundary := s.nextDistance(posX, posY, dir)
		if dist == 0 {
			if boundary {
				timeLeft = 0
				break
			}
		}

		if timeLeft < dist {
			posX += dirDx[dir] * int(timeLeft)
			posY += dirDy[dir] * int(timeLeft)
			timeLeft = 0
			break
		}

		move := int(dist)
		posX += dirDx[dir] * move
		posY += dirDy[dir] * move
		elapsed += dist
		timeLeft -= dist
		lastArrowID = -1

		if boundary {
			timeLeft = 0
			break
		}
	}
	return posX, posY
}

type fastScanner struct {
	reader *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{reader: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.reader.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = fs.reader.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.reader.ReadByte()
	}
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.reader.ReadByte()
	}
	return sign * val
}

func (fs *fastScanner) nextInt64() int64 {
	sign := int64(1)
	var val int64
	c, err := fs.reader.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = fs.reader.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.reader.ReadByte()
	}
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.reader.ReadByte()
	}
	return sign * val
}

func (fs *fastScanner) nextString() string {
	c, err := fs.reader.ReadByte()
	for err == nil && c <= ' ' {
		c, err = fs.reader.ReadByte()
	}
	buf := make([]byte, 0, 16)
	for err == nil && c > ' ' {
		buf = append(buf, c)
		c, err = fs.reader.ReadByte()
	}
	return string(buf)
}

func directionFromChar(ch byte) int {
	switch ch {
	case 'L':
		return dirLeft
	case 'U':
		return dirUp
	case 'D':
		return dirDown
	default:
		return dirRight
	}
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := fs.nextInt()
	b := fs.nextInt()

	s := newSolver(b)
	for i := 0; i < n; i++ {
		x0 := fs.nextInt()
		y0 := fs.nextInt()
		x1 := fs.nextInt()
		y1 := fs.nextInt()
		s.addArrow(x0, y0, x1, y1)
	}
	s.build()

	q := fs.nextInt()
	for i := 0; i < q; i++ {
		x := fs.nextInt()
		y := fs.nextInt()
		dirStr := fs.nextString()
		t := fs.nextInt64()
		dir := directionFromChar(dirStr[0])
		resX, resY := s.query(x, y, dir, t)
		fmt.Fprintf(out, "%d %d\n", resX, resY)
	}
}

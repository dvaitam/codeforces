package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int64 {
	sign := int64(1)
	var val int64
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *FastScanner) NextString() string {
	c, _ := fs.r.ReadByte()
	for c == ' ' || c == '\n' || c == '\r' || c == '\t' {
		c, _ = fs.r.ReadByte()
	}
	buf := []byte{c}
	for {
		c, err := fs.r.ReadByte()
		if err != nil || c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			break
		}
		buf = append(buf, c)
	}
	return string(buf)
}

const (
	DIR_UP = iota
	DIR_DOWN
	DIR_LEFT
	DIR_RIGHT
)

var dx = [4]int64{0, 0, -1, 1}
var dy = [4]int64{1, -1, 0, 0}

func dirDelta(dir int) (int64, int64) {
	return dx[dir], dy[dir]
}

type Arrow struct {
	id     int
	x0, y0 int64
	x1, y1 int64
	vert   bool
	dir    int
	exitX  int64
	exitY  int64
	xConst int64
	yConst int64
	yLow   int64
	yHigh  int64
	xLow   int64
	xHigh  int64
}

type HorQuery struct {
	y   int64
	x   int64
	dir int
	idx int
}

type VerQuery struct {
	x   int64
	y   int64
	dir int
	idx int
}

type NextEvent struct {
	typ      int
	boundary int64
	arrow    int
}

type VertSeg struct {
	id    int
	x     int64
	yLow  int64
	yHigh int64
}

type HorSeg struct {
	id    int
	y     int64
	xLow  int64
	xHigh int64
}

type Edge struct {
	dist  int64
	first int64
	arrow int
	to    int
}

type TreapNode struct {
	key      int64
	val      int
	priority uint32
	left     *TreapNode
	right    *TreapNode
}

type Treap struct {
	root *TreapNode
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func rotateLeft(x *TreapNode) *TreapNode {
	y := x.right
	x.right = y.left
	y.left = x
	return y
}

func rotateRight(x *TreapNode) *TreapNode {
	y := x.left
	x.left = y.right
	y.right = x
	return y
}

func treapInsert(node *TreapNode, key int64, val int) *TreapNode {
	if node == nil {
		return &TreapNode{key: key, val: val, priority: rng.Uint32()}
	}
	if key < node.key {
		node.left = treapInsert(node.left, key, val)
		if node.left.priority < node.priority {
			node = rotateRight(node)
		}
	} else if key > node.key {
		node.right = treapInsert(node.right, key, val)
		if node.right.priority < node.priority {
			node = rotateLeft(node)
		}
	} else {
		node.val = val
	}
	return node
}

func treapMerge(a, b *TreapNode) *TreapNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		a.right = treapMerge(a.right, b)
		return a
	}
	b.left = treapMerge(a, b.left)
	return b
}

func treapDelete(node *TreapNode, key int64) *TreapNode {
	if node == nil {
		return nil
	}
	if key < node.key {
		node.left = treapDelete(node.left, key)
	} else if key > node.key {
		node.right = treapDelete(node.right, key)
	} else {
		node = treapMerge(node.left, node.right)
	}
	return node
}

func (t *Treap) Insert(key int64, val int) {
	t.root = treapInsert(t.root, key, val)
}

func (t *Treap) Delete(key int64) {
	t.root = treapDelete(t.root, key)
}

func (t *Treap) Successor(key int64) (int64, int) {
	node := t.root
	var best *TreapNode
	for node != nil {
		if node.key > key {
			best = node
			node = node.left
		} else {
			node = node.right
		}
	}
	if best == nil {
		return 0, -1
	}
	return best.key, best.val
}

func (t *Treap) Predecessor(key int64) (int64, int) {
	node := t.root
	var best *TreapNode
	for node != nil {
		if node.key < key {
			best = node
			node = node.right
		} else {
			node = node.left
		}
	}
	if best == nil {
		return 0, -1
	}
	return best.key, best.val
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func processHorizontal(verticals []VertSeg, requests []HorQuery, boundary int64) []NextEvent {
	res := make([]NextEvent, len(requests))
	if len(requests) == 0 {
		return res
	}
	yset := make(map[int64]struct{})
	for _, seg := range verticals {
		yset[seg.yLow] = struct{}{}
		yset[seg.yHigh] = struct{}{}
	}
	for _, q := range requests {
		yset[q.y] = struct{}{}
	}
	ys := make([]int64, 0, len(yset))
	for y := range yset {
		ys = append(ys, y)
	}
	// sort ys
	for i := 1; i < len(ys); i++ {
		key := ys[i]
		j := i - 1
		for j >= 0 && ys[j] > key {
			ys[j+1] = ys[j]
			j--
		}
		ys[j+1] = key
	}
	index := make(map[int64]int, len(ys))
	for i, y := range ys {
		index[y] = i
	}
	add := make([][]int, len(ys))
	rem := make([][]int, len(ys))
	qs := make([][]HorQuery, len(ys))
	for idx, seg := range verticals {
		l := index[seg.yLow]
		r := index[seg.yHigh]
		add[l] = append(add[l], idx)
		rem[r] = append(rem[r], idx)
	}
	for _, q := range requests {
		pos := index[q.y]
		qs[pos] = append(qs[pos], q)
	}
	treap := &Treap{}
	for i := 0; i < len(ys); i++ {
		for _, id := range add[i] {
			treap.Insert(verticals[id].x, verticals[id].id)
		}
		for _, q := range qs[i] {
			if q.dir == DIR_RIGHT {
				_, arrow := treap.Successor(q.x)
				if arrow == -1 {
					res[q.idx] = NextEvent{typ: 0, boundary: boundary}
				} else {
					res[q.idx] = NextEvent{typ: 1, arrow: arrow}
				}
			} else {
				_, arrow := treap.Predecessor(q.x)
				if arrow == -1 {
					res[q.idx] = NextEvent{typ: 0, boundary: 0}
				} else {
					res[q.idx] = NextEvent{typ: 1, arrow: arrow}
				}
			}
		}
		for _, id := range rem[i] {
			treap.Delete(verticals[id].x)
		}
	}
	return res
}

func processVertical(horizontals []HorSeg, requests []VerQuery, boundary int64) []NextEvent {
	res := make([]NextEvent, len(requests))
	if len(requests) == 0 {
		return res
	}
	xset := make(map[int64]struct{})
	for _, seg := range horizontals {
		xset[seg.xLow] = struct{}{}
		xset[seg.xHigh] = struct{}{}
	}
	for _, q := range requests {
		xset[q.x] = struct{}{}
	}
	xs := make([]int64, 0, len(xset))
	for x := range xset {
		xs = append(xs, x)
	}
	for i := 1; i < len(xs); i++ {
		key := xs[i]
		j := i - 1
		for j >= 0 && xs[j] > key {
			xs[j+1] = xs[j]
			j--
		}
		xs[j+1] = key
	}
	index := make(map[int64]int, len(xs))
	for i, x := range xs {
		index[x] = i
	}
	add := make([][]int, len(xs))
	rem := make([][]int, len(xs))
	qs := make([][]VerQuery, len(xs))
	for idx, seg := range horizontals {
		l := index[seg.xLow]
		r := index[seg.xHigh]
		add[l] = append(add[l], idx)
		rem[r] = append(rem[r], idx)
	}
	for _, q := range requests {
		pos := index[q.x]
		qs[pos] = append(qs[pos], q)
	}
	treap := &Treap{}
	for i := 0; i < len(xs); i++ {
		for _, id := range add[i] {
			treap.Insert(horizontals[id].y, horizontals[id].id)
		}
		for _, q := range qs[i] {
			if q.dir == DIR_UP {
				_, arrow := treap.Successor(q.y)
				if arrow == -1 {
					res[q.idx] = NextEvent{typ: 0, boundary: boundary}
				} else {
					res[q.idx] = NextEvent{typ: 1, arrow: arrow}
				}
			} else {
				_, arrow := treap.Predecessor(q.y)
				if arrow == -1 {
					res[q.idx] = NextEvent{typ: 0, boundary: 0}
				} else {
					res[q.idx] = NextEvent{typ: 1, arrow: arrow}
				}
			}
		}
		for _, id := range rem[i] {
			treap.Delete(horizontals[id].y)
		}
	}
	return res
}

func buildEdge(arr Arrow, res NextEvent, arrows []Arrow) Edge {
	startX := arr.exitX
	startY := arr.exitY
	dir := arr.dir
	if res.typ == 0 {
		var dist int64
		if dir == DIR_LEFT || dir == DIR_RIGHT {
			dist = abs64(res.boundary - startX)
		} else {
			dist = abs64(res.boundary - startY)
		}
		return Edge{dist: dist, first: dist, arrow: -1, to: -1}
	}
	next := res.arrow
	var firstLen int64
	if dir == DIR_LEFT || dir == DIR_RIGHT {
		target := arrows[next].xConst
		firstLen = abs64(target - startX)
	} else {
		target := arrows[next].yConst
		firstLen = abs64(target - startY)
	}
	entryX := startX + dx[dir]*firstLen
	entryY := startY + dy[dir]*firstLen
	exitX := arrows[next].exitX
	exitY := arrows[next].exitY
	arrowLen := abs64(exitX-entryX) + abs64(exitY-entryY)
	total := firstLen + arrowLen
	return Edge{dist: total, first: firstLen, arrow: next, to: next}
}

type Query struct {
	x, y int64
	dir  int
	t    int64
}

func moveAlong(x, y int64, dir int, dist int64) (int64, int64) {
	return x + dx[dir]*dist, y + dy[dir]*dist
}

func processInitial(x, y int64, dir int, t int64, res NextEvent, arrows []Arrow) (bool, int64, int64, int, int64) {
	if res.typ == 0 {
		var dist int64
		if dir == DIR_LEFT || dir == DIR_RIGHT {
			dist = abs64(res.boundary - x)
		} else {
			dist = abs64(res.boundary - y)
		}
		if t <= dist {
			nx, ny := moveAlong(x, y, dir, t)
			return true, nx, ny, -1, 0
		}
		nx, ny := moveAlong(x, y, dir, dist)
		return true, nx, ny, -1, 0
	}
	next := res.arrow
	var firstLen int64
	if dir == DIR_LEFT || dir == DIR_RIGHT {
		target := arrows[next].xConst
		firstLen = abs64(target - x)
	} else {
		target := arrows[next].yConst
		firstLen = abs64(target - y)
	}
	entryX := x + dx[dir]*firstLen
	entryY := y + dy[dir]*firstLen
	exitX := arrows[next].exitX
	exitY := arrows[next].exitY
	arrowLen := abs64(exitX-entryX) + abs64(exitY-entryY)
	total := firstLen + arrowLen
	if t <= firstLen {
		nx, ny := moveAlong(x, y, dir, t)
		return true, nx, ny, -1, 0
	}
	if t <= total {
		extra := t - firstLen
		nx := entryX + dx[arrows[next].dir]*extra
		ny := entryY + dy[arrows[next].dir]*extra
		return true, nx, ny, -1, 0
	}
	remaining := t - total
	return false, arrows[next].exitX, arrows[next].exitY, next, remaining
}

func simulateState(state int, remaining int64, arrows []Arrow, edges []Edge, jump [][]int, jumpDist [][]int64) (int64, int64) {
	curState := state
	curX := arrows[curState].exitX
	curY := arrows[curState].exitY
	for {
		if remaining == 0 {
			return curX, curY
		}
		for k := len(jump) - 1; k >= 0; k-- {
			nxt := jump[k][curState]
			if nxt != -1 && jumpDist[k][curState] <= remaining {
				remaining -= jumpDist[k][curState]
				curState = nxt
				curX = arrows[curState].exitX
				curY = arrows[curState].exitY
			}
		}
		edge := edges[curState]
		if edge.arrow == -1 {
			if remaining >= edge.dist {
				nx := curX + dx[arrows[curState].dir]*edge.dist
				ny := curY + dy[arrows[curState].dir]*edge.dist
				return nx, ny
			}
			nx := curX + dx[arrows[curState].dir]*remaining
			ny := curY + dy[arrows[curState].dir]*remaining
			return nx, ny
		}
		if remaining >= edge.dist {
			remaining -= edge.dist
			curState = edge.to
			curX = arrows[curState].exitX
			curY = arrows[curState].exitY
			continue
		}
		if remaining <= edge.first {
			nx := curX + dx[arrows[curState].dir]*remaining
			ny := curY + dy[arrows[curState].dir]*remaining
			return nx, ny
		}
		entryX := curX + dx[arrows[curState].dir]*edge.first
		entryY := curY + dy[arrows[curState].dir]*edge.first
		extra := remaining - edge.first
		nx := entryX + dx[arrows[edge.arrow].dir]*extra
		ny := entryY + dy[arrows[edge.arrow].dir]*extra
		return nx, ny
	}
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := int(fs.NextInt())
	b := fs.NextInt()
	arrows := make([]Arrow, n)
	var verticals []VertSeg
	var horizontals []HorSeg
	stateHorIdx := make([]int, n)
	stateVerIdx := make([]int, n)
	for i := range stateHorIdx {
		stateHorIdx[i] = -1
	}
	for i := range stateVerIdx {
		stateVerIdx[i] = -1
	}
	for i := 0; i < n; i++ {
		x0 := fs.NextInt()
		y0 := fs.NextInt()
		x1 := fs.NextInt()
		y1 := fs.NextInt()
		arr := Arrow{id: i, x0: x0, y0: y0, x1: x1, y1: y1}
		if x0 == x1 {
			arr.vert = true
			arr.xConst = x0
			if y1 > y0 {
				arr.dir = DIR_UP
			} else {
				arr.dir = DIR_DOWN
			}
			arr.exitX = x1
			arr.exitY = y1
			if y0 < y1 {
				arr.yLow = y0
				arr.yHigh = y1
			} else {
				arr.yLow = y1
				arr.yHigh = y0
			}
			verticals = append(verticals, VertSeg{id: i, x: arr.xConst, yLow: arr.yLow, yHigh: arr.yHigh})
		} else {
			arr.vert = false
			arr.yConst = y0
			if x1 > x0 {
				arr.dir = DIR_RIGHT
			} else {
				arr.dir = DIR_LEFT
			}
			arr.exitX = x1
			arr.exitY = y1
			if x0 < x1 {
				arr.xLow = x0
				arr.xHigh = x1
			} else {
				arr.xLow = x1
				arr.xHigh = x0
			}
			horizontals = append(horizontals, HorSeg{id: i, y: arr.yConst, xLow: arr.xLow, xHigh: arr.xHigh})
		}
		arrows[i] = arr
	}

	var horRequests []HorQuery
	var verRequests []VerQuery
	for _, arr := range arrows {
		if arr.vert {
			idx := len(verRequests)
			stateVerIdx[arr.id] = idx
			verRequests = append(verRequests, VerQuery{
				x:   arr.exitX,
				y:   arr.exitY,
				dir: arr.dir,
				idx: idx,
			})
		} else {
			idx := len(horRequests)
			stateHorIdx[arr.id] = idx
			horRequests = append(horRequests, HorQuery{
				y:   arr.exitY,
				x:   arr.exitX,
				dir: arr.dir,
				idx: idx,
			})
		}
	}

	q := int(fs.NextInt())
	queries := make([]Query, q)
	queryHorIdx := make([]int, q)
	queryVerIdx := make([]int, q)
	for i := 0; i < q; i++ {
		x := fs.NextInt()
		y := fs.NextInt()
		dirStr := fs.NextString()
		t := fs.NextInt()
		dir := 0
		switch dirStr[0] {
		case 'U':
			dir = DIR_UP
		case 'D':
			dir = DIR_DOWN
		case 'L':
			dir = DIR_LEFT
		case 'R':
			dir = DIR_RIGHT
		}
		queries[i] = Query{x: x, y: y, dir: dir, t: t}
		if dir == DIR_LEFT || dir == DIR_RIGHT {
			idx := len(horRequests)
			queryHorIdx[i] = idx
			horRequests = append(horRequests, HorQuery{
				y:   y,
				x:   x,
				dir: dir,
				idx: idx,
			})
		} else {
			idx := len(verRequests)
			queryVerIdx[i] = idx
			verRequests = append(verRequests, VerQuery{
				x:   x,
				y:   y,
				dir: dir,
				idx: idx,
			})
		}
	}

	horResults := processHorizontal(verticals, horRequests, b)
	verResults := processVertical(horizontals, verRequests, b)

	edges := make([]Edge, n)
	for i, arr := range arrows {
		if arr.vert {
			idx := stateVerIdx[i]
			if idx >= 0 {
				edges[i] = buildEdge(arr, verResults[idx], arrows)
			}
		} else {
			idx := stateHorIdx[i]
			if idx >= 0 {
				edges[i] = buildEdge(arr, horResults[idx], arrows)
			}
		}
	}

	const LOG = 60
	jump := make([][]int, LOG)
	jumpDist := make([][]int64, LOG)
	for k := 0; k < LOG; k++ {
		jump[k] = make([]int, n)
		jumpDist[k] = make([]int64, n)
		for i := 0; i < n; i++ {
			jump[k][i] = -1
		}
	}
	for i := 0; i < n; i++ {
		jump[0][i] = edges[i].to
		jumpDist[0][i] = edges[i].dist
	}
	for k := 1; k < LOG; k++ {
		for i := 0; i < n; i++ {
			if jump[k-1][i] == -1 {
				jump[k][i] = -1
				jumpDist[k][i] = jumpDist[k-1][i]
			} else {
				jump[k][i] = jump[k-1][jump[k-1][i]]
				jumpDist[k][i] = jumpDist[k-1][i] + jumpDist[k-1][jump[k-1][i]]
			}
		}
	}

	for i := 0; i < q; i++ {
		query := queries[i]
		var res NextEvent
		if query.dir == DIR_LEFT || query.dir == DIR_RIGHT {
			res = horResults[queryHorIdx[i]]
		} else {
			res = verResults[queryVerIdx[i]]
		}
		done, fx, fy, nextState, remaining := processInitial(query.x, query.y, query.dir, query.t, res, arrows)
		if done {
			fmt.Fprintf(out, "%d %d\n", fx, fy)
			continue
		}
		if nextState == -1 {
			fmt.Fprintf(out, "%d %d\n", fx, fy)
			continue
		}
		ansX, ansY := simulateState(nextState, remaining, arrows, edges, jump, jumpDist)
		fmt.Fprintf(out, "%d %d\n", ansX, ansY)
	}
}

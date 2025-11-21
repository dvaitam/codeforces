package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf = int64(1 << 60)
const negInf = -inf

type interval struct {
	l int64
	r int64
}

type component struct {
	L [2]int64
	R [2]int64
}

type compIntervals struct {
	A interval
	B interval
}

type segment struct {
	start int64
	yL    int64
	yR    int64
}

type event struct {
	start int64
	comp  int
	yL    int64
	yR    int64
}

type segTree struct {
	size int
	max  []int64
	min  []int64
}

func newSegTree(n int) *segTree {
	size := 1
	for size < n {
		size <<= 1
	}
	maxArr := make([]int64, size<<1)
	minArr := make([]int64, size<<1)
	for i := range maxArr {
		maxArr[i] = negInf
		minArr[i] = inf
	}
	return &segTree{
		size: size,
		max:  maxArr,
		min:  minArr,
	}
}

func (st *segTree) update(pos int, maxVal, minVal int64) {
	idx := pos + st.size
	st.max[idx] = maxVal
	st.min[idx] = minVal
	idx >>= 1
	for idx > 0 {
		left := idx << 1
		right := left | 1
		if st.max[left] > st.max[right] {
			st.max[idx] = st.max[left]
		} else {
			st.max[idx] = st.max[right]
		}
		if st.min[left] < st.min[right] {
			st.min[idx] = st.min[left]
		} else {
			st.min[idx] = st.min[right]
		}
		idx >>= 1
	}
}

func (st *segTree) getMax() int64 {
	return st.max[1]
}

func (st *segTree) getMin() int64 {
	return st.min[1]
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func buildSegments(a, b interval, limit int64) []segment {
	points := []int64{0, limit + 1, a.l, a.r + 1, b.l, b.r + 1}
	sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
	uniq := make([]int64, 0, len(points))
	for _, val := range points {
		if len(uniq) == 0 || uniq[len(uniq)-1] != val {
			uniq = append(uniq, val)
		}
	}
	segs := make([]segment, 0, len(uniq)-1)
	for i := 0; i < len(uniq)-1; i++ {
		l := uniq[i]
		r := uniq[i+1] - 1
		if r < 0 || l > limit {
			continue
		}
		if l < 0 {
			l = 0
		}
		if r > limit {
			r = limit
		}
		if l > r {
			continue
		}
		x := l
		inA := a.l <= x && x <= a.r
		inB := b.l <= x && x <= b.r
		var yL, yR int64
		if inA && inB {
			if a.l < b.l {
				yL = a.l
			} else {
				yL = b.l
			}
			if a.r > b.r {
				yR = a.r
			} else {
				yR = b.r
			}
		} else if inA {
			yL, yR = b.l, b.r
		} else if inB {
			yL, yR = a.l, a.r
		} else {
			yL = limit + 1
			yR = -1
		}
		segs = append(segs, segment{start: l, yL: yL, yR: yR})
	}
	return segs
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var tInt, TInt int64
	if _, err := fmt.Fscan(in, &tInt, &TInt); err != nil {
		return
	}
	var n, m int
	fmt.Fscan(in, &n, &m)
	ls := make([]int64, n)
	rs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ls[i], &rs[i])
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	compID := make([]int, n)
	components := make([]component, 0)
	queue := make([]int, 0)

	for i := 0; i < n; i++ {
		if color[i] != -1 {
			continue
		}
		comp := component{
			L: [2]int64{0, 0},
			R: [2]int64{inf, inf},
		}
		color[i] = 0
		compID[i] = len(components)
		queue = queue[:0]
		queue = append(queue, i)
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			c := color[v]
			if ls[v] > comp.L[c] {
				comp.L[c] = ls[v]
			}
			if rs[v] < comp.R[c] {
				comp.R[c] = rs[v]
			}
			for _, to := range adj[v] {
				if color[to] == -1 {
					color[to] = c ^ 1
					compID[to] = len(components)
					queue = append(queue, to)
				} else if color[to] == c {
					fmt.Println("IMPOSSIBLE")
					return
				}
			}
		}
		if comp.L[0] > comp.R[0] || comp.L[1] > comp.R[1] {
			fmt.Println("IMPOSSIBLE")
			return
		}
		components = append(components, comp)
	}

	numComp := len(components)
	compIntervalsArr := make([]compIntervals, numComp)
	for i, comp := range components {
		aL := max64(comp.L[0], 0)
		aR := min64(comp.R[0], TInt)
		bL := max64(comp.L[1], 0)
		bR := min64(comp.R[1], TInt)
		if aL > aR || bL > bR {
			fmt.Println("IMPOSSIBLE")
			return
		}
		compIntervalsArr[i] = compIntervals{
			A: interval{l: aL, r: aR},
			B: interval{l: bL, r: bR},
		}
	}

	events := make([]event, 0, numComp*5)
	for idx, iv := range compIntervalsArr {
		segs := buildSegments(iv.A, iv.B, TInt)
		sort.Slice(segs, func(i, j int) bool { return segs[i].start < segs[j].start })
		for _, seg := range segs {
			events = append(events, event{
				start: seg.start,
				comp:  idx,
				yL:    seg.yL,
				yR:    seg.yR,
			})
		}
	}

	if len(events) == 0 {
		fmt.Println("IMPOSSIBLE")
		return
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].start == events[j].start {
			return events[i].comp < events[j].comp
		}
		return events[i].start < events[j].start
	})

	starts := make([]int64, 0)
	for i := 0; i < len(events); {
		start := events[i].start
		starts = append(starts, start)
		for i < len(events) && events[i].start == start {
			i++
		}
	}
	if starts[0] != 0 {
		fmt.Println("IMPOSSIBLE")
		return
	}
	if starts[len(starts)-1] != TInt+1 {
		starts = append(starts, TInt+1)
	}

	st := newSegTree(numComp)
	invalidL := TInt + 1
	invalidR := int64(-1)
	for i := 0; i < numComp; i++ {
		st.update(i, invalidL, invalidR)
	}

	eventIdx := 0
	found := false
	var finalX, finalY int64

	for idx := 0; idx < len(starts)-1; idx++ {
		start := starts[idx]
		for eventIdx < len(events) && events[eventIdx].start == start {
			ev := events[eventIdx]
			st.update(ev.comp, ev.yL, ev.yR)
			eventIdx++
		}
		if start > TInt {
			break
		}
		nextStart := starts[idx+1]
		if nextStart > TInt+1 {
			nextStart = TInt + 1
		}
		xL := start
		xR := nextStart - 1
		if xR > TInt {
			xR = TInt
		}
		if xL > xR {
			continue
		}
		maxL := st.getMax()
		minR := st.getMin()
		if maxL > minR {
			continue
		}
		yL := maxL
		yR := minR
		loX := max64(xL, tInt-yR)
		hiX := min64(xR, TInt-yL)
		if loX > hiX {
			continue
		}
		xCand := loX
		yCandLow := max64(yL, tInt-xCand)
		yCandHigh := min64(yR, TInt-xCand)
		if yCandLow > yCandHigh {
			xCand = hiX
			yCandLow = max64(yL, tInt-xCand)
			yCandHigh = min64(yR, TInt-xCand)
			if yCandLow > yCandHigh {
				continue
			}
		}
		finalX = xCand
		finalY = yCandLow
		found = true
		break
	}

	if !found {
		fmt.Println("IMPOSSIBLE")
		return
	}

	total := finalX + finalY
	if total < tInt || total > TInt {
		fmt.Println("IMPOSSIBLE")
		return
	}

	compFlip := make([]int, numComp)
	for idx, iv := range compIntervalsArr {
		inA := iv.A.l <= finalX && finalX <= iv.A.r
		inB := iv.B.l <= finalX && finalX <= iv.B.r
		switch {
		case inA && !inB:
			compFlip[idx] = 0
			if finalY < iv.B.l || finalY > iv.B.r {
				fmt.Println("IMPOSSIBLE")
				return
			}
		case !inA && inB:
			compFlip[idx] = 1
			if finalY < iv.A.l || finalY > iv.A.r {
				fmt.Println("IMPOSSIBLE")
				return
			}
		case inA && inB:
			if iv.B.l <= finalY && finalY <= iv.B.r {
				compFlip[idx] = 0
			} else if iv.A.l <= finalY && finalY <= iv.A.r {
				compFlip[idx] = 1
			} else {
				fmt.Println("IMPOSSIBLE")
				return
			}
		default:
			fmt.Println("IMPOSSIBLE")
			return
		}
	}

	assign := make([]byte, n)
	for i := 0; i < n; i++ {
		comp := compID[i]
		group := color[i] ^ compFlip[comp]
		assign[i] = byte('1' + group)
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, "POSSIBLE")
	fmt.Fprintf(out, "%d %d\n", finalX, finalY)
	fmt.Fprintln(out, string(assign))
	out.Flush()
}

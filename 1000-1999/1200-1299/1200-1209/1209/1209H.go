package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
}

func nextInt() int {
	scanner.Scan()
	v, _ := strconv.Atoi(scanner.Text())
	return v
}

func nextFloat() float64 {
	scanner.Scan()
	v, _ := strconv.ParseFloat(scanner.Text(), 64)
	return v
}

type SegTree struct {
	minVal []float64
	lazy   []float64
}

func NewSegTree(n int) *SegTree {
	return &SegTree{
		minVal: make([]float64, 4*n+1),
		lazy:   make([]float64, 4*n+1),
	}
}

func (st *SegTree) Build(node, l, r int, E []float64) {
	if l == r {
		st.minVal[node] = E[l]
		return
	}
	mid := (l + r) / 2
	st.Build(2*node, l, mid, E)
	st.Build(2*node+1, mid+1, r, E)
	st.minVal[node] = math.Min(st.minVal[2*node], st.minVal[2*node+1])
}

func (st *SegTree) pushDown(node int) {
	if st.lazy[node] != 0 {
		st.lazy[2*node] += st.lazy[node]
		st.minVal[2*node] += st.lazy[node]

		st.lazy[2*node+1] += st.lazy[node]
		st.minVal[2*node+1] += st.lazy[node]

		st.lazy[node] = 0
	}
}

func (st *SegTree) Add(node, l, r, ql, qr int, val float64) {
	if ql <= l && r <= qr {
		st.lazy[node] += val
		st.minVal[node] += val
		return
	}
	st.pushDown(node)
	mid := (l + r) / 2
	if ql <= mid {
		st.Add(2*node, l, mid, ql, qr, val)
	}
	if qr > mid {
		st.Add(2*node+1, mid+1, r, ql, qr, val)
	}
	st.minVal[node] = math.Min(st.minVal[2*node], st.minVal[2*node+1])
}

func (st *SegTree) QueryMin(node, l, r, ql, qr int) float64 {
	if ql <= l && r <= qr {
		return st.minVal[node]
	}
	st.pushDown(node)
	mid := (l + r) / 2
	res := math.MaxFloat64
	if ql <= mid {
		res = math.Min(res, st.QueryMin(2*node, l, mid, ql, qr))
	}
	if qr > mid {
		res = math.Min(res, st.QueryMin(2*node+1, mid+1, r, ql, qr))
	}
	return res
}

type Segment struct {
	d     float64
	s     float64
	x     float64
	max_x float64
}

func main() {
	if !scanner.Scan() {
		return
	}
	N, _ := strconv.Atoi(scanner.Text())
	L := nextInt()

	var segments []Segment
	curr := 0

	for i := 0; i < N; i++ {
		X := nextInt()
		Y := nextInt()
		S := nextFloat()
		if X > curr {
			segments = append(segments, Segment{
				d: float64(X - curr),
				s: 0,
			})
		}
		segments = append(segments, Segment{
			d: float64(Y - X),
			s: S,
		})
		curr = Y
	}
	if curr < L {
		segments = append(segments, Segment{
			d: float64(L - curr),
			s: 0,
		})
	}

	M := len(segments)
	E := make([]float64, M+1)
	for i := 0; i < M; i++ {
		d := segments[i].d
		s := segments[i].s
		if s > 0 {
			segments[i].x = -d / s
		} else {
			segments[i].x = 0
		}
		segments[i].max_x = d / (s + 2.0)
		E[i+1] = E[i] - segments[i].x
	}

	st := NewSegTree(M)
	st.Build(1, 1, M, E)

	sortedIndices := make([]int, M)
	for i := range sortedIndices {
		sortedIndices[i] = i
	}
	sort.Slice(sortedIndices, func(i, j int) bool {
		return segments[sortedIndices[i]].s < segments[sortedIndices[j]].s
	})

	for _, idx := range sortedIndices {
		seg := &segments[idx]
		allowance := st.QueryMin(1, 1, M, idx+1, M)
		max_inc := seg.max_x - seg.x
		if max_inc > 0 && allowance > 0 {
			inc := math.Min(allowance, max_inc)
			seg.x += inc
			st.Add(1, 1, M, idx+1, M, -inc)
		}
	}

	totalTime := 0.0
	for i := 0; i < M; i++ {
		totalTime += (segments[i].d - segments[i].x) / (segments[i].s + 1.0)
	}
	fmt.Printf("%.15f\n", totalTime)
}

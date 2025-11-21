package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type State struct {
	lamHigh float64
	lamLow  float64
	L       float64
	U       float64
	w       float64
	w2      float64
}

type Event struct {
	lambda float64
	typ    int // +1 activation, -1 deactivation
	state  *State
}

type Interval struct {
	xStart float64
	xEnd   float64
	C0     float64
	C1     float64
	D0     float64
}

const eps = 1e-12

func invSqrt(x float64) float64 {
	if math.IsInf(x, 1) {
		return 0
	}
	if x <= 0 {
		return math.Inf(1)
	}
	return 1.0 / math.Sqrt(x)
}

func recordInterval(intervals *[]Interval, lamHi, lamLo, C0, C1, D0 float64) {
	if !(lamHi > lamLo+eps) || C1 <= eps {
		return
	}
	xStart := C0 + C1*invSqrt(lamHi)
	xEnd := C0 + C1*invSqrt(lamLo)
	if xEnd <= xStart+1e-15 {
		return
	}
	*intervals = append(*intervals, Interval{
		xStart: xStart,
		xEnd:   xEnd,
		C0:     C0,
		C1:     C1,
		D0:     D0,
	})
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t, q int
	if _, err := fmt.Fscan(in, &t, &q); err != nil {
		return
	}

	states := make([]*State, 0)
	events := make([]Event, 0)

	costAtZero := 0.0
	totalLimit := 0.0

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]float64, n)
		for j := 0; j < n; j++ {
			var v int
			fmt.Fscan(in, &v)
			a[j] = float64(v)
		}
		pPrime := make([]float64, n)
		sumP := 0.0
		for j := 0; j < n; j++ {
			var v int
			fmt.Fscan(in, &v)
			pPrime[j] = float64(v)
			sumP += pPrime[j]
		}
		p := make([]float64, n)
		for j := 0; j < n; j++ {
			p[j] = pPrime[j] / sumP
		}

		type pair struct {
			th float64
			a  float64
			p  float64
		}
		items := make([]pair, n)
		sumA := 0.0
		maxTh := 0.0
		for j := 0; j < n; j++ {
			sumA += a[j]
		}
		for j := 0; j < n; j++ {
			th := a[j] / p[j]
			items[j] = pair{th: th, a: a[j], p: p[j]}
			if th > maxTh {
				maxTh = th
			}
		}
		totalLimit += math.Max(0, maxTh-sumA)

		sort.Slice(items, func(i, j int) bool {
			if math.Abs(items[i].th-items[j].th) < 1e-15 {
				return items[i].p > items[j].p
			}
			return items[i].th > items[j].th
		})

		prefixA := make([]float64, n+1)
		prefixP := make([]float64, n+1)
		for j := 0; j < n; j++ {
			prefixA[j+1] = prefixA[j] + items[j].a
			prefixP[j+1] = prefixP[j] + items[j].p
		}

		m0 := 0
		for j := 0; j < n; j++ {
			if items[j].th > sumA+1e-15 {
				m0 = j + 1
			} else {
				break
			}
		}

		if sumA == 0 {
			continue
		}

		if m0 > 0 {
			Sa := prefixA[m0]
			Sp := prefixP[m0]
			costAtZero += 2*Sa/sumA - 2*Sp
		}

		for k := m0; k >= 1; k-- {
			var L float64
			if k == m0 {
				L = sumA
			} else {
				L = items[k].th
			}
			U := items[k-1].th
			if U <= L+1e-15 {
				continue
			}
			Sa := prefixA[k]
			w2 := 2 * Sa
			w := math.Sqrt(w2)
			state := &State{
				lamHigh: w2 / (L * L),
				lamLow:  w2 / (U * U),
				L:       L,
				U:       U,
				w:       w,
				w2:      w2,
			}
			states = append(states, state)
			events = append(events, Event{lambda: state.lamHigh, typ: +1, state: state})
			events = append(events, Event{lambda: state.lamLow, typ: -1, state: state})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if math.Abs(events[i].lambda-events[j].lambda) < 1e-18 {
			return events[i].typ > events[j].typ
		}
		return events[i].lambda > events[j].lambda
	})

	intervals := make([]Interval, 0)

	C0 := 0.0
	C1 := 0.0
	D0 := costAtZero
	lamPrev := math.Inf(1)

	idx := 0
	for {
		lamCurr := 0.0
		if idx < len(events) {
			lamCurr = events[idx].lambda
		}
		recordInterval(&intervals, lamPrev, lamCurr, C0, C1, D0)
		if idx == len(events) {
			break
		}
		for idx < len(events) && math.Abs(events[idx].lambda-lamCurr) < 1e-18 {
			st := events[idx].state
			if events[idx].typ == +1 {
				C1 += st.w
				C0 -= st.L
				D0 -= st.w2 / st.L
			} else {
				C1 -= st.w
				C0 += st.U
				D0 += st.w2 / st.U
			}
			idx++
		}
		lamPrev = lamCurr
	}

	xLimit := totalLimit

	xs := make([]float64, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &xs[i])
	}

	for _, x := range xs {
		if x <= 0 {
			fmt.Fprintf(out, "%.12f\n", costAtZero)
			continue
		}
		if x >= xLimit-1e-9 {
			fmt.Fprintf(out, "0\n")
			continue
		}
		pos := sort.Search(len(intervals), func(i int) bool {
			return x < intervals[i].xEnd-1e-12
		})
		if pos == len(intervals) {
			fmt.Fprintf(out, "0\n")
			continue
		}
		inter := intervals[pos]
		denom := x - inter.C0
		if denom <= 0 {
			denom = 1e-18
		}
		cost := inter.D0 + inter.C1*inter.C1/denom
		fmt.Fprintf(out, "%.12f\n", cost)
	}
}

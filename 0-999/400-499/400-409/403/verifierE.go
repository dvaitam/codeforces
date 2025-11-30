package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"5",
	"1 1 2 4",
	"1 2 2 3",
	"2",
	"",
	"6",
	"1 2 1 1 5",
	"1 1 2 1 1",
	"3",
	"",
	"5",
	"1 2 2 3",
	"1 2 2 3",
	"1",
	"",
	"6",
	"1 1 3 4 1",
	"1 2 1 3 1",
	"2",
	"",
	"6",
	"1 1 1 4 1",
	"1 2 3 4 1",
	"3",
	"",
	"6",
	"1 1 3 3 5",
	"1 2 2 1 5",
	"4",
	"",
	"4",
	"1 2 1",
	"1 1 1",
	"3",
	"",
	"4",
	"1 1 1",
	"1 1 1",
	"1",
	"",
	"6",
	"1 2 3 2 2",
	"1 2 2 4 3",
	"1",
	"",
	"4",
	"1 2 3",
	"1 1 1",
	"1",
	"",
	"4",
	"1 1 2",
	"1 2 2",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 1 1 2 5",
	"1 2 1 3 1",
	"1",
	"",
	"6",
	"1 1 1 1 4",
	"1 1 3 1 5",
	"4",
	"",
	"6",
	"1 2 1 2 1",
	"1 2 2 2 1",
	"5",
	"",
	"5",
	"1 1 3 4",
	"1 2 2 4",
	"2",
	"",
	"3",
	"1 1",
	"1 2",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 2 3 3",
	"1 2 1 1",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"2",
	"",
	"6",
	"1 2 2 4 1",
	"1 1 3 3 2",
	"2",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"6",
	"1 1 1 4 1",
	"1 1 2 4 5",
	"5",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 1 1 2 3",
	"1 1 1 4 4",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 1",
	"1 2 1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 3",
	"1 2 2",
	"3",
	"",
	"3",
	"1 2",
	"1 1",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"1",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"2",
	"",
	"4",
	"1 1 2",
	"1 2 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 1",
	"3",
	"",
	"5",
	"1 2 2 1",
	"1 1 3 4",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"6",
	"1 2 2 2 4",
	"1 2 1 3 5",
	"2",
	"",
	"5",
	"1 2 1 1",
	"1 2 3 2",
	"1",
	"",
	"5",
	"1 2 1 1",
	"1 1 1 2",
	"4",
	"",
	"5",
	"1 1 1 4",
	"1 2 3 4",
	"4",
	"",
	"4",
	"1 2 3",
	"1 1 1",
	"2",
	"",
	"4",
	"1 1 3",
	"1 2 3",
	"1",
	"",
	"5",
	"1 2 1 1",
	"1 1 1 2",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 3",
	"2",
	"",
	"4",
	"1 1 3",
	"1 1 3",
	"2",
	"",
	"3",
	"1 2",
	"1 2",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 2 1 4 3",
	"1 2 2 4 1",
	"2",
	"",
	"3",
	"1 2",
	"1 1",
	"1",
	"",
	"5",
	"1 1 2 1",
	"1 2 2 1",
	"4",
	"",
	"5",
	"1 1 1 4",
	"1 1 1 2",
	"3",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 1 3 4",
	"1 1 3 3",
	"2",
	"",
	"5",
	"1 1 3 2",
	"1 2 2 3",
	"2",
	"",
	"5",
	"1 1 1 1",
	"1 2 1 2",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 2",
	"1",
	"",
	"6",
	"1 1 1 4 1",
	"1 2 3 3 2",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"2",
	"",
	"6",
	"1 1 3 1 1",
	"1 1 2 2 2",
	"2",
	"",
	"4",
	"1 1 1",
	"1 2 2",
	"3",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 1 3",
	"1 2 3",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 2 2 1 4",
	"1 2 1 3 3",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"4",
	"1 2 2",
	"1 2 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 1",
	"2",
	"",
	"5",
	"1 2 1 4",
	"1 1 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 2 2 3",
	"1 1 1 1",
	"3",
	"",
	"4",
	"1 1 3",
	"1 1 1",
	"2",
	"",
	"4",
	"1 1 3",
	"1 1 2",
	"3",
	"",
	"6",
	"1 2 2 4 3",
	"1 1 1 1 4",
	"4",
	"",
	"6",
	"1 1 3 3 3",
	"1 2 1 4 4",
	"5",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 1 3 1",
	"1 2 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"5",
	"1 2 1 3",
	"1 2 1 4",
	"2",
	"",
	"6",
	"1 1 1 2 4",
	"1 2 2 4 4",
	"2",
	"",
}

type testcase struct {
	n     int
	a     []int
	b     []int
	idx   int
	input string
}

func parseInts(line string) []int {
	if strings.TrimSpace(line) == "" {
		return nil
	}
	fields := strings.Fields(line)
	res := make([]int, len(fields))
	for i, f := range fields {
		val, _ := strconv.Atoi(f)
		res[i] = val
	}
	return res
}

func parseCases() []testcase {
	var cases []testcase
	for pos := 0; pos < len(rawTestcases); {
		for pos < len(rawTestcases) && strings.TrimSpace(rawTestcases[pos]) == "" {
			pos++
		}
		if pos >= len(rawTestcases) {
			break
		}
		n, _ := strconv.Atoi(strings.TrimSpace(rawTestcases[pos]))
		pos++
		if pos+2 >= len(rawTestcases) {
			break
		}
		lineA := rawTestcases[pos]
		lineB := rawTestcases[pos+1]
		lineIdx := rawTestcases[pos+2]
		pos += 3

		aParents := parseInts(lineA)
		bParents := parseInts(lineB)
		idxVal, _ := strconv.Atoi(strings.TrimSpace(lineIdx))

		a := make([]int, n+1)
		b := make([]int, n+1)
		for i := 0; i < len(aParents) && i+2 <= n; i++ {
			a[i+2] = aParents[i]
		}
		for i := 0; i < len(bParents) && i+2 <= n; i++ {
			b[i+2] = bParents[i]
		}

		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		fmt.Fprintln(&sb, strings.TrimSpace(lineA))
		fmt.Fprintln(&sb, strings.TrimSpace(lineB))
		fmt.Fprintln(&sb, strings.TrimSpace(lineIdx))

		cases = append(cases, testcase{
			n:     n,
			a:     a,
			b:     b,
			idx:   idxVal,
			input: sb.String(),
		})
	}
	return cases
}

// --- begin solution from 403E.go adapted as function ---
type SegMax struct {
	n    int
	tree []int
}

func NewSegMax(n int) *SegMax {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegMax{n: size, tree: make([]int, 2*size)}
}
func (st *SegMax) Update(pos, val int) {
	i := pos + st.n - 1
	st.tree[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if st.tree[2*i] > st.tree[2*i+1] {
			st.tree[i] = st.tree[2*i]
		} else {
			st.tree[i] = st.tree[2*i+1]
		}
	}
}
func (st *SegMax) QueryFirst(l, r, thr int) int { return st.query(1, 1, st.n, l, r, thr) }
func (st *SegMax) query(idx, lo, hi, l, r, thr int) int {
	if lo > r || hi < l || st.tree[idx] <= thr {
		return 0
	}
	if lo == hi {
		return lo
	}
	mid := (lo + hi) >> 1
	if res := st.query(2*idx, lo, mid, l, r, thr); res != 0 {
		return res
	}
	return st.query(2*idx+1, mid+1, hi, l, r, thr)
}

type SegMin struct {
	n    int
	tree []int
}

func NewSegMin(n int) *SegMin {
	size := 1
	for size < n {
		size <<= 1
	}
	inf := int(1e9)
	tree := make([]int, 2*size)
	for i := range tree {
		tree[i] = inf
	}
	return &SegMin{n: size, tree: tree}
}
func (st *SegMin) Update(pos, val int) {
	i := pos + st.n - 1
	st.tree[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if st.tree[2*i] < st.tree[2*i+1] {
			st.tree[i] = st.tree[2*i]
		} else {
			st.tree[i] = st.tree[2*i+1]
		}
	}
}
func (st *SegMin) QueryFirst(l, r, thr int) int { return st.query(1, 1, st.n, l, r, thr) }
func (st *SegMin) query(idx, lo, hi, l, r, thr int) int {
	if lo > r || hi < l || st.tree[idx] >= thr {
		return 0
	}
	if lo == hi {
		return lo
	}
	mid := (lo + hi) >> 1
	if res := st.query(2*idx, lo, mid, l, r, thr); res != 0 {
		return res
	}
	return st.query(2*idx+1, mid+1, hi, l, r, thr)
}

func solve(n int, a, b []int, idx0 int) string {
	var out bytes.Buffer
	blueAdj := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		blueAdj[a[i]] = append(blueAdj[a[i]], i)
	}
	blueTin := make([]int, n+1)
	blueTout := make([]int, n+1)
	timer := 0
	var dfsBlue func(int)
	dfsBlue = func(u int) {
		timer++
		blueTin[u] = timer
		for _, v := range blueAdj[u] {
			dfsBlue(v)
		}
		blueTout[u] = timer
	}
	dfsBlue(1)
	redAdj := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		redAdj[b[i]] = append(redAdj[b[i]], i)
	}
	redTin := make([]int, n+1)
	redTout := make([]int, n+1)
	timer = 0
	var dfsRed func(int)
	dfsRed = func(u int) {
		timer++
		redTin[u] = timer
		for _, v := range redAdj[u] {
			dfsRed(v)
		}
		redTout[u] = timer
	}
	dfsRed(1)
	aaR := make([]int, n)
	bbR := make([]int, n)
	listARA := make([][]int, n+2)
	listARB := make([][]int, n+2)
	for ei := 1; ei < n; ei++ {
		u := b[ei+1]
		v := ei + 1
		t1 := blueTin[u]
		t2 := blueTin[v]
		if t1 < t2 {
			aaR[ei] = t1
			bbR[ei] = t2
		} else {
			aaR[ei] = t2
			bbR[ei] = t1
		}
		listARA[aaR[ei]] = append(listARA[aaR[ei]], ei)
		listARB[bbR[ei]] = append(listARB[bbR[ei]], ei)
	}
	for i := 1; i <= n; i++ {
		if len(listARA[i]) > 1 {
			sort.Slice(listARA[i], func(x, y int) bool { return bbR[listARA[i][x]] > bbR[listARA[i][y]] })
		}
		if len(listARB[i]) > 1 {
			sort.Slice(listARB[i], func(x, y int) bool { return aaR[listARB[i][x]] < aaR[listARB[i][y]] })
		}
	}
	segRA := NewSegMax(n)
	segRB := NewSegMin(n)
	ptrRA := make([]int, n+2)
	ptrRB := make([]int, n+2)
	const INF = int(1e9)
	for i := 1; i <= n; i++ {
		if ptrRA[i] < len(listARA[i]) {
			eid := listARA[i][ptrRA[i]]
			segRA.Update(i, bbR[eid])
		}
		if ptrRB[i] < len(listARB[i]) {
			eid := listARB[i][ptrRB[i]]
			segRB.Update(i, aaR[eid])
		}
	}
	aaB := make([]int, n)
	bbB := make([]int, n)
	listBAA := make([][]int, n+2)
	listBAB := make([][]int, n+2)
	for ei := 1; ei < n; ei++ {
		u := a[ei+1]
		v := ei + 1
		t1 := redTin[u]
		t2 := redTin[v]
		if t1 < t2 {
			aaB[ei] = t1
			bbB[ei] = t2
		} else {
			aaB[ei] = t2
			bbB[ei] = t1
		}
		listBAA[aaB[ei]] = append(listBAA[aaB[ei]], ei)
		listBAB[bbB[ei]] = append(listBAB[bbB[ei]], ei)
	}
	for i := 1; i <= n; i++ {
		if len(listBAA[i]) > 1 {
			sort.Slice(listBAA[i], func(x, y int) bool { return bbB[listBAA[i][x]] > bbB[listBAA[i][y]] })
		}
		if len(listBAB[i]) > 1 {
			sort.Slice(listBAB[i], func(x, y int) bool { return aaB[listBAB[i][x]] < aaB[listBAB[i][y]] })
		}
	}
	segBA := NewSegMax(n)
	segBB := NewSegMin(n)
	ptrBA := make([]int, n+2)
	ptrBB := make([]int, n+2)
	for i := 1; i <= n; i++ {
		if ptrBA[i] < len(listBAA[i]) {
			eid := listBAA[i][ptrBA[i]]
			segBA.Update(i, bbB[eid])
		}
		if ptrBB[i] < len(listBAB[i]) {
			eid := listBAB[i][ptrBB[i]]
			segBB.Update(i, aaB[eid])
		}
	}
	deletedR := make([]bool, n)
	deletedB := make([]bool, n)
	type Stage struct {
		color string
		edges []int
	}
	var stages []Stage
	blueQ := []int{idx0}
	for len(blueQ) > 0 {
		sort.Ints(blueQ)
		stages = append(stages, Stage{"Blue", append([]int(nil), blueQ...)})
		redQ := make([]int, 0)
		for _, be := range blueQ {
			v := be + 1
			l, r := blueTin[v], blueTout[v]
			for {
				bb := segRB.QueryFirst(l, r, l)
				if bb == 0 {
					break
				}
				for ptrRB[bb] < len(listARB[bb]) {
					e := listARB[bb][ptrRB[bb]]
					if deletedR[e] || aaR[e] >= l {
						ptrRB[bb]++
						continue
					}
					deletedR[e] = true
					redQ = append(redQ, e)
					ptrRB[bb]++
					if ptrRB[bb] < len(listARB[bb]) {
						segRB.Update(bb, aaR[listARB[bb][ptrRB[bb]]])
					} else {
						segRB.Update(bb, INF)
					}
					aa0 := aaR[e]
					for ptrRA[aa0] < len(listARA[aa0]) && deletedR[listARA[aa0][ptrRA[aa0]]] {
						ptrRA[aa0]++
					}
					if ptrRA[aa0] < len(listARA[aa0]) {
						segRA.Update(aa0, bbR[listARA[aa0][ptrRA[aa0]]])
					} else {
						segRA.Update(aa0, 0)
					}
					break
				}
			}
			for {
				aa := segRA.QueryFirst(l, r, r)
				if aa == 0 {
					break
				}
				for ptrRA[aa] < len(listARA[aa]) {
					e := listARA[aa][ptrRA[aa]]
					if deletedR[e] || bbR[e] <= r {
						ptrRA[aa]++
						continue
					}
					deletedR[e] = true
					redQ = append(redQ, e)
					ptrRA[aa]++
					if ptrRA[aa] < len(listARA[aa]) {
						segRA.Update(aa, bbR[listARA[aa][ptrRA[aa]]])
					} else {
						segRA.Update(aa, 0)
					}
					bb0 := bbR[e]
					for ptrRB[bb0] < len(listARB[bb0]) && deletedR[listARB[bb0][ptrRB[bb0]]] {
						ptrRB[bb0]++
					}
					if ptrRB[bb0] < len(listARB[bb0]) {
						segRB.Update(bb0, aaR[listARB[bb0][ptrRB[bb0]]])
					} else {
						segRB.Update(bb0, INF)
					}
					break
				}
			}
		}
		if len(redQ) == 0 {
			break
		}
		sort.Ints(redQ)
		stages = append(stages, Stage{"Red", append([]int(nil), redQ...)})
		blueQ = make([]int, 0)
		for _, re := range redQ {
			v := re + 1
			l, r := redTin[v], redTout[v]
			for {
				bb := segBB.QueryFirst(l, r, l)
				if bb == 0 {
					break
				}
				for ptrBB[bb] < len(listBAB[bb]) {
					e := listBAB[bb][ptrBB[bb]]
					if deletedB[e] || aaB[e] >= l {
						ptrBB[bb]++
						continue
					}
					deletedB[e] = true
					blueQ = append(blueQ, e)
					ptrBB[bb]++
					if ptrBB[bb] < len(listBAB[bb]) {
						segBB.Update(bb, aaB[listBAB[bb][ptrBB[bb]]])
					} else {
						segBB.Update(bb, INF)
					}
					aa0 := aaB[e]
					for ptrBA[aa0] < len(listBAA[aa0]) && deletedB[listBAA[aa0][ptrBA[aa0]]] {
						ptrBA[aa0]++
					}
					if ptrBA[aa0] < len(listBAA[aa0]) {
						segBA.Update(aa0, bbB[listBAA[aa0][ptrBA[aa0]]])
					} else {
						segBA.Update(aa0, 0)
					}
					break
				}
			}
			for {
				aa := segBA.QueryFirst(l, r, r)
				if aa == 0 {
					break
				}
				for ptrBA[aa] < len(listBAA[aa]) {
					e := listBAA[aa][ptrBA[aa]]
					if deletedB[e] || bbB[e] <= r {
						ptrBA[aa]++
						continue
					}
					deletedB[e] = true
					blueQ = append(blueQ, e)
					ptrBA[aa]++
					if ptrBA[aa] < len(listBAA[aa]) {
						segBA.Update(aa, bbB[listBAA[aa][ptrBA[aa]]])
					} else {
						segBA.Update(aa, 0)
					}
					bb0 := bbB[e]
					for ptrBB[bb0] < len(listBAB[bb0]) && deletedB[listBAB[bb0][ptrBB[bb0]]] {
						ptrBB[bb0]++
					}
					if ptrBB[bb0] < len(listBAB[bb0]) {
						segBB.Update(bb0, aaB[listBAB[bb0][ptrBB[bb0]]])
					} else {
						segBB.Update(bb0, INF)
					}
					break
				}
			}
		}
	}
	for _, stg := range stages {
		fmt.Fprintln(&out, stg.color)
		for i, e := range stg.edges {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(&out, e)
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

// --- end adapted solution ---

func runSolution(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for i, tc := range cases {
		expected := solve(tc.n, tc.a, tc.b, tc.idx)
		got, err := runSolution(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d mismatch:\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

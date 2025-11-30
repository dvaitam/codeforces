package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesD.txt.
const embeddedTestcasesD = `5 5 2 4 5 3 6 3 0 3 10 8 11 1 6 16 1 0 2 7 8 9
3 4 0 1 7 1 7 4 10 13 7 3 12 5 2 5 6
3 4 6 7 4 8 13 8 7 8 8 9 18 4 0 7 10
4 1 6 12 7 0 4 9 10
3 5 4 8 8 8 14 6 4 8 1 9 15 3 9 18 10 2 5 9
3 1 9 10 9 2 5 9
3 3 5 13 6 9 16 2 0 3 9 0 2 7
3 1 7 12 9 1 3 8
4 3 2 11 2 1 4 9 1 10 10 6 8 9 10
3 3 4 5 6 4 9 8 8 17 5 3 5 6
2 1 8 11 10 6 9
3 4 0 9 6 10 20 3 0 6 8 2 6 10 0 5 7
2 5 3 8 2 6 8 9 7 11 7 8 10 3 2 10 1 6 10
3 3 6 12 9 5 7 4 0 8 0 3 4 6
4 4 10 17 0 9 17 5 9 12 9 4 10 0 2 6 7 8
1 1 9 15 5 0
1 2 1 10 7 0 6 0 5
4 2 10 15 6 10 13 9 0 2 4 6
2 2 2 10 10 0 9 0 6 8
2 3 9 11 1 8 11 4 3 8 5 4 10
5 4 2 10 8 2 3 10 9 12 10 8 9 5 1 3 6 7 9
4 5 2 8 10 2 10 8 0 9 1 8 14 0 1 3 6 5 7 9 10
4 5 5 7 2 5 6 2 2 3 5 9 13 0 6 7 4 0 2 5 6
1 4 0 8 5 9 19 4 10 15 9 7 14 2 0
4 3 3 10 1 5 7 1 0 6 0 0 2 5 6
4 5 7 15 1 0 9 6 4 5 10 8 10 1 5 11 1 0 2 4 7
1 1 6 12 2 8
2 2 2 5 10 10 14 9 0 5
4 4 0 4 3 10 15 5 2 6 5 3 6 6 2 5 6 7
5 1 2 12 0 0 2 6 8 10
1 3 8 9 0 0 2 9 9 12 2 10
4 1 6 13 9 2 3 5 10
3 5 3 12 6 1 4 6 9 15 1 6 13 3 7 14 3 3 6 7
4 5 1 6 4 8 14 8 0 10 9 7 11 4 0 10 5 0 1 6 8
2 4 0 7 6 6 8 7 9 17 2 2 8 7 2 6
5 3 8 10 5 5 8 10 5 13 10 0 1 3 4 8
5 3 4 11 10 0 5 8 6 7 10 2 3 4 6 10
2 2 1 11 5 2 3 6 6 9
4 1 10 12 6 2 8 9 10
2 2 2 6 0 8 11 7 5 9
3 3 10 12 6 4 7 5 10 19 5 2 6 8
5 3 3 10 5 6 14 8 4 11 6 0 1 2 4 9
5 5 1 5 9 10 20 8 1 6 9 2 9 1 0 1 1 1 5 7 8 10
4 3 9 14 10 7 11 2 8 17 1 0 2 8 9
5 4 9 13 7 10 17 4 0 10 2 6 9 7 0 1 5 6 9
5 4 5 9 8 7 8 7 9 11 4 8 16 8 1 4 6 7 10
2 5 5 8 4 9 12 7 1 3 7 6 16 6 8 10 4 3 7
1 3 2 8 1 2 3 10 2 12 8 3
1 1 6 15 9 7
1 4 8 14 5 1 2 3 3 11 4 4 8 0 7
3 5 5 7 1 4 14 6 3 9 6 2 6 4 3 11 10 4 5 6
5 2 1 8 5 8 16 3 2 5 6 9 10
3 4 4 6 7 4 6 7 2 8 3 2 8 7 1 3 6
4 4 8 16 9 9 13 10 6 15 4 7 11 5 0 1 7 8
3 4 3 10 0 9 10 6 1 6 3 5 8 1 0 2 5
2 1 0 7 8 0 2
1 5 9 13 1 7 11 0 8 18 6 1 10 2 3 7 6 6
4 1 6 10 6 0 4 9 10
5 3 5 11 10 7 10 9 8 10 4 0 1 4 5 10
2 5 10 13 6 3 13 10 5 9 6 8 17 1 8 10 7 1 8
4 4 2 10 8 5 8 6 4 11 1 9 18 5 3 5 7 10
4 4 0 8 0 8 15 7 3 10 3 4 12 7 2 3 4 7
3 4 9 12 8 10 12 3 4 13 10 10 12 0 0 2 5
3 5 2 10 6 10 14 6 7 11 2 6 7 5 10 19 3 2 5 8
5 5 9 16 2 10 20 7 10 20 3 9 17 0 3 11 9 0 3 5 7 8
1 2 3 11 10 9 13 10 10
1 5 6 11 10 5 13 10 3 6 5 8 14 7 6 14 10 5
4 1 2 6 1 2 3 6 7
3 1 0 7 3 1 6 9
5 2 10 11 10 2 11 10 0 1 4 6 7
5 4 0 4 6 9 15 4 1 2 3 9 10 9 2 5 7 8 10
1 5 7 12 3 6 13 8 5 9 9 4 6 0 2 8 10 9
5 4 1 10 7 6 12 6 0 2 6 9 12 2 0 1 3 6 10
3 3 6 8 5 4 10 9 10 14 8 0 5 8
5 5 0 3 7 0 5 5 2 8 1 9 18 6 1 6 9 1 2 3 8 10
1 5 3 9 1 8 12 2 10 14 6 3 13 7 5 8 4 5
5 2 3 11 5 2 7 8 0 1 2 5 10
3 1 1 11 0 1 2 4
3 2 8 17 5 7 9 6 6 8 10
1 1 9 17 5 2
3 2 10 17 10 4 5 0 0 1 4
3 3 8 13 0 3 12 2 4 14 0 4 6 10
5 2 4 10 8 7 15 1 0 1 2 5 6
4 4 0 4 1 4 6 2 2 10 0 5 14 6 1 3 7 8
3 3 6 13 7 0 3 8 5 9 3 2 3 4
2 2 8 13 8 0 3 3 2 6
2 5 2 12 10 2 6 8 4 5 7 3 4 6 2 9 8 3 6
4 5 8 16 1 2 4 8 0 2 0 3 10 7 10 15 6 1 4 5 6
2 5 3 7 0 0 9 1 7 11 7 10 13 9 10 19 4 0 9
2 3 10 16 5 4 13 10 4 7 0 3 7
5 4 5 13 4 7 10 7 4 10 5 10 20 6 4 5 6 9 10
5 2 3 10 3 4 13 2 2 7 8 9 10
3 4 4 12 10 8 9 3 4 12 1 0 4 9 0 1 4
2 2 9 12 10 5 15 4 0 10
1 5 2 7 3 6 14 5 5 7 9 6 8 0 2 3 9 8
3 3 3 10 3 5 8 2 7 10 9 6 8 9
5 2 0 1 6 9 10 10 1 2 3 8 10
3 3 10 20 3 0 6 2 5 10 7 5 7 9
3 1 7 10 1 4 8 10
4 5 7 14 10 2 12 6 2 4 9 1 4 7 4 6 10 1 4 6 9`

type wall struct {
	l int64
	r int64
	t int64
}

type minHeap []int64

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveCase(n int, walls []wall, qs []int64) []int64 {
	ys := make([]int64, 0, 2*len(walls))
	for _, w := range walls {
		ys = append(ys, w.l, w.r)
	}
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	uniq := ys[:0]
	for i, v := range ys {
		if i == 0 || v != ys[i-1] {
			uniq = append(uniq, v)
		}
	}
	ys = uniq
	idx := func(y int64) int {
		return sort.Search(len(ys), func(i int) bool { return ys[i] >= y })
	}
	adds := make([][]int64, len(ys))
	rems := make([][]int64, len(ys))
	for _, w := range walls {
		li := idx(w.l)
		ri := idx(w.r)
		adds[li] = append(adds[li], w.t)
		rems[ri] = append(rems[ri], w.t)
	}
	active := &minHeap{}
	heap.Init(active)
	removeCnt := make(map[int64]int)
	type seg struct{ A, B, L int64 }
	segs := make([]seg, 0, len(ys))
	for k := 0; k < len(ys); k++ {
		for _, t := range rems[k] {
			removeCnt[t]++
		}
		for _, t := range adds[k] {
			heap.Push(active, t)
		}
		if k+1 < len(ys) {
			for active.Len() > 0 {
				top := (*active)[0]
				if removeCnt[top] > 0 {
					heap.Pop(active)
					removeCnt[top]--
				} else {
					break
				}
			}
			if active.Len() > 0 {
				M := (*active)[0]
				y0 := ys[k]
				y1 := ys[k+1]
				L := y1 - y0
				A := M - y1
				B := A + L
				segs = append(segs, seg{A, B, L})
			}
		}
	}
	As := make([]int64, len(segs))
	Bs := make([]struct{ B, A, L int64 }, len(segs))
	for i, s := range segs {
		As[i] = s.A
		Bs[i] = struct{ B, A, L int64 }{s.B, s.A, s.L}
	}
	sort.Slice(As, func(i, j int) bool { return As[i] < As[j] })
	sort.Slice(Bs, func(i, j int) bool { return Bs[i].B < Bs[j].B })

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return qs[order[i]] < qs[order[j]] })

	ans := make([]int64, n)
	pa, pb := 0, 0
	var cntA, cntB int64
	var sumA, sumAB, sumFullL int64
	for _, oi := range order {
		q := qs[oi]
		for pa < len(As) && As[pa] <= q {
			sumA += As[pa]
			cntA++
			pa++
		}
		for pb < len(Bs) && Bs[pb].B <= q {
			sumAB += Bs[pb].A
			sumFullL += Bs[pb].L
			cntB++
			pb++
		}
		ans[oi] = sumFullL + (cntA-cntB)*q - (sumA - sumAB)
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func parseLine(line string) (int, []wall, []int64, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0, nil, nil, fmt.Errorf("invalid line")
	}
	n64, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, nil, nil, err
	}
	m64, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, nil, nil, err
	}
	n := int(n64)
	m := int(m64)
	expected := 2 + 3*m + n
	if len(parts) != expected {
		return 0, nil, nil, fmt.Errorf("expected %d values got %d", expected, len(parts))
	}
	walls := make([]wall, m)
	pos := 2
	for i := 0; i < m; i++ {
		l, _ := strconv.ParseInt(parts[pos], 10, 64)
		r, _ := strconv.ParseInt(parts[pos+1], 10, 64)
		t, _ := strconv.ParseInt(parts[pos+2], 10, 64)
		walls[i] = wall{l: l, r: r, t: t}
		pos += 3
	}
	qs := make([]int64, n)
	for i := 0; i < n; i++ {
		q, _ := strconv.ParseInt(parts[pos+i], 10, 64)
		qs[i] = q
	}
	return n, walls, qs, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesD), "\n")
	for idx, line := range lines {
		n, walls, qs, err := parseLine(strings.TrimSpace(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expectedVals := solveCase(n, walls, qs)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, len(walls))
		for _, w := range walls {
			fmt.Fprintf(&input, "%d %d %d\n", w.l, w.r, w.t)
		}
		for _, q := range qs {
			fmt.Fprintf(&input, "%d ", q)
		}
		input.WriteByte('\n')

		var want strings.Builder
		for i, v := range expectedVals {
			if i > 0 {
				want.WriteByte('\n')
			}
			want.WriteString(strconv.FormatInt(v, 10))
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want.String()) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, strings.TrimSpace(want.String()), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(strings.Split(strings.TrimSpace(embeddedTestcasesD), "\n")))
}

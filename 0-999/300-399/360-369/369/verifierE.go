package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Segment struct{ l, r int }

func NewBIT(n int) *BIT { return &BIT{n: n, t: make([]int, n+2)} }

type BIT struct {
	n int
	t []int
}

func (b *BIT) Add(i, v int) {
	for i <= b.n {
		b.t[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	if i > b.n {
		i = b.n
	}
	s := 0
	for i > 0 {
		s += b.t[i]
		i -= i & -i
	}
	return s
}

func computeE(n, m int, segs []Segment, queries [][]int) []int {
	lArr := make([]int, n)
	rArr := make([]int, n)
	for i, s := range segs {
		lArr[i] = s.l
		rArr[i] = s.r
	}
	sort.Ints(lArr)
	sort.Ints(rArr)
	ansL := make([]int, m)
	ansR := make([]int, m)
	gapSum := make([]int, m)
	type Gap struct{ x, y, idx int }
	var gaps []Gap
	for qi, pts := range queries {
		sort.Ints(pts)
		p1 := pts[0]
		ansL[qi] = sort.SearchInts(rArr, p1)
		pk := pts[len(pts)-1]
		idxR := sort.Search(len(lArr), func(i int) bool { return lArr[i] > pk })
		ansR[qi] = n - idxR
		for i := 0; i+1 < len(pts); i++ {
			gaps = append(gaps, Gap{pts[i], pts[i+1], qi})
		}
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].l > segs[j].l })
	sort.Slice(gaps, func(i, j int) bool { return gaps[i].x > gaps[j].x })
	maxC := 1000
	bit := NewBIT(maxC)
	si := 0
	for _, g := range gaps {
		for si < n && segs[si].l > g.x {
			r := segs[si].r
			if r >= 1 && r <= maxC {
				bit.Add(r, 1)
			}
			si++
		}
		lim := g.y - 1
		if lim > 0 {
			gapSum[g.idx] += bit.Sum(lim)
		}
	}
	res := make([]int, m)
	for i := 0; i < m; i++ {
		res[i] = n - ansL[i] - ansR[i] - gapSum[i]
	}
	return res
}

type caseE struct {
	n, m    int
	segs    []Segment
	queries [][]int
	input   string
	expect  []int
}

func generateCase(rng *rand.Rand) caseE {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	segs := make([]Segment, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		l := rng.Intn(20) + 1
		r := l + rng.Intn(10)
		segs[i] = Segment{l, r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	queries := make([][]int, m)
	for qi := 0; qi < m; qi++ {
		cnt := rng.Intn(3) + 1
		pts := make([]int, cnt)
		for i := 0; i < cnt; i++ {
			pts[i] = rng.Intn(25)
		}
		queries[qi] = pts
		sb.WriteString(fmt.Sprintf("%d", cnt))
		for _, p := range pts {
			sb.WriteString(fmt.Sprintf(" %d", p))
		}
		sb.WriteByte('\n')
	}
	expect := computeE(n, m, segs, queries)
	return caseE{n, m, segs, queries, sb.String(), expect}
}

func runCase(bin string, c caseE) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(c.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := strings.NewReader(out.String())
	for i := 0; i < c.m; i++ {
		var got int
		if _, err := fmt.Fscan(reader, &got); err != nil {
			return fmt.Errorf("failed to read answer %d: %v\n%s", i+1, err, out.String())
		}
		if got != c.expect[i] {
			return fmt.Errorf("answer %d expected %d got %d", i+1, c.expect[i], got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

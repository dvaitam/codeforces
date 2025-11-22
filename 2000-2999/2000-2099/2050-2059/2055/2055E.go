package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type scanner struct {
	r *bufio.Reader
}

func newScanner() *scanner {
	return &scanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (s *scanner) nextInt64() int64 {
	sign, val := int64(1), int64(0)
	c, _ := s.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = s.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = s.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = s.r.ReadByte()
	}
	return sign * val
}

type pair struct {
	a   int64
	b   int64
	idx int
}

func buildSequence(a, b []int64) ([]pair, []int) {
	n := len(a)
	good := make([]pair, 0, n)
	bad := make([]pair, 0, n)
	for i := 0; i < n; i++ {
		if b[i] >= a[i] {
			good = append(good, pair{a[i], b[i], i})
		} else {
			bad = append(bad, pair{a[i], b[i], i})
		}
	}
	sort.Slice(good, func(i, j int) bool {
		if good[i].a == good[j].a {
			return good[i].idx < good[j].idx
		}
		return good[i].a < good[j].a
	})
	sort.Slice(bad, func(i, j int) bool {
		if bad[i].b == bad[j].b {
			return bad[i].idx < bad[j].idx
		}
		return bad[i].b > bad[j].b
	})
	seq := append(good, bad...)
	pos := make([]int, n)
	for i, p := range seq {
		pos[p.idx] = i
	}
	return seq, pos
}

func solveCase(a, b []int64) int64 {
	n := len(a)
	seq, pos := buildSequence(a, b)
	m := len(seq)

	prefCur := make([]int64, m+1)
	prefExtra := make([]int64, m+1)
	for i := 0; i < m; i++ {
		cur := prefCur[i]
		extra := prefExtra[i]
		if cur < seq[i].a {
			extra += seq[i].a - cur
			cur = seq[i].a
		}
		cur = cur - seq[i].a + seq[i].b
		prefCur[i+1] = cur
		prefExtra[i+1] = extra
	}

	need := make([]int64, m+1)
	delta := make([]int64, m+1)
	for i := m - 1; i >= 0; i-- {
		delta[i] = delta[i+1] + (seq[i].b - seq[i].a)
		val := need[i+1] + seq[i].a - seq[i].b
		if val < seq[i].a {
			val = seq[i].a
		}
		need[i] = val
	}

	total := int64(0)
	for _, v := range a {
		total += v
	}

	best := int64(-1)
	for k := 0; k < n; k++ {
		p := pos[k]
		cur := prefCur[p]
		extra := prefExtra[p]

		req := need[p+1]
		add := int64(0)
		if cur < req {
			add = req - cur
			cur = req
		}
		cur += delta[p+1]
		totalExtra := extra + add
		if cur < a[k]+totalExtra {
			continue
		}
		moves := total + totalExtra
		if best == -1 || moves < best {
			best = moves
		}
	}
	return best
}

func main() {
	in := newScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int(in.nextInt64())
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = in.nextInt64()
			b[i] = in.nextInt64()
		}
		ans := solveCase(a, b)
		fmt.Fprintln(out, ans)
	}
}

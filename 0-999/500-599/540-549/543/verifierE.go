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

type SegTree struct {
	n    int
	minv []int
	lazy []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	minv := make([]int, size*2)
	lazy := make([]int, size*2)
	return &SegTree{n: n, minv: minv, lazy: lazy}
}

func (st *SegTree) apply(p, v int) {
	st.minv[p] += v
	st.lazy[p] += v
}

func (st *SegTree) push(p int) {
	if st.lazy[p] != 0 {
		st.apply(p*2, st.lazy[p])
		st.apply(p*2+1, st.lazy[p])
		st.lazy[p] = 0
	}
}

func (st *SegTree) pull(p int) {
	if st.minv[p*2] < st.minv[p*2+1] {
		st.minv[p] = st.minv[p*2]
	} else {
		st.minv[p] = st.minv[p*2+1]
	}
}

func (st *SegTree) updateRange(p, l, r, ql, qr, v int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(p, v)
		return
	}
	st.push(p)
	m := (l + r) >> 1
	st.updateRange(p*2, l, m, ql, qr, v)
	st.updateRange(p*2+1, m+1, r, ql, qr, v)
	st.pull(p)
}

func (st *SegTree) queryMin(p, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return 1 << 60
	}
	if ql <= l && r <= qr {
		return st.minv[p]
	}
	st.push(p)
	m := (l + r) >> 1
	a := st.queryMin(p*2, l, m, ql, qr)
	b := st.queryMin(p*2+1, m+1, r, ql, qr)
	if a < b {
		return a
	}
	return b
}

func solveE(n, m int, a []int, qs [][3]int) []int {
	W := n - m + 1
	if W < 1 {
		W = 1
	}
	type P struct{ v, idx int }
	ps := make([]P, n)
	for j := 0; j < n; j++ {
		ps[j] = P{v: a[j], idx: j}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v })
	st := NewSegTree(W)
	p := 0
	curQ := 0
	prevAns := 0
	res := make([]int, len(qs))
	for i, q := range qs {
		qi := q[2] ^ prevAns
		if qi > curQ {
			for p < n && ps[p].v < qi {
				j := ps[p].idx
				l := j - m + 1
				if l < 0 {
					l = 0
				}
				r := j
				if r >= W {
					r = W - 1
				}
				if l <= r {
					st.updateRange(1, 0, st.n-1, l, r, 1)
				}
				p++
			}
		} else if qi < curQ {
			for p > 0 && ps[p-1].v >= qi {
				p--
				j := ps[p].idx
				l := j - m + 1
				if l < 0 {
					l = 0
				}
				r := j
				if r >= W {
					r = W - 1
				}
				if l <= r {
					st.updateRange(1, 0, st.n-1, l, r, -1)
				}
			}
		}
		curQ = qi
		li := q[0]
		ri := q[1]
		if ri >= W {
			ri = W - 1
		}
		ans := st.queryMin(1, 0, st.n-1, li, ri)
		if ans < 0 {
			ans = 0
		}
		res[i] = ans
		prevAns = ans
	}
	return res
}

func genTestE() (string, string) {
	n := rand.Intn(5) + 1
	m := rand.Intn(n) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rand.Intn(10)
	}
	s := rand.Intn(5) + 1
	qs := make([][3]int, s)
	for i := 0; i < s; i++ {
		l := rand.Intn(n)
		r := rand.Intn(n)
		if l > r {
			l, r = r, l
		}
		x := rand.Intn(10)
		qs[i] = [3]int{l, r, x}
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	for i := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", a[i])
	}
	input += "\n"
	input += fmt.Sprintf("%d\n", s)
	for i := 0; i < s; i++ {
		input += fmt.Sprintf("%d %d %d\n", qs[i][0]+1, qs[i][1]+1, qs[i][2])
	}
	res := solveE(n, m, a, qs)
	expected := ""
	for i, v := range res {
		if i > 0 {
			expected += " "
		}
		expected += fmt.Sprintf("%d", v)
	}
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestE()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %s\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

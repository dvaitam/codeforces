package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const INF = math.MaxInt64 / 4

type segTree struct {
	n    int
	min  []int64
	idx  []int
	lazy []int64
}

func newSegTree(a []int64) *segTree {
	n := len(a) - 1
	size := 4 * (n + 1)
	st := &segTree{n: n, min: make([]int64, size), idx: make([]int, size), lazy: make([]int64, size)}
	var build func(node, l, r int)
	build = func(node, l, r int) {
		if l == r {
			st.min[node] = a[l]
			st.idx[node] = l
		} else {
			m := (l + r) >> 1
			build(node<<1, l, m)
			build(node<<1|1, m+1, r)
			if st.min[node<<1] <= st.min[node<<1|1] {
				st.min[node] = st.min[node<<1]
				st.idx[node] = st.idx[node<<1]
			} else {
				st.min[node] = st.min[node<<1|1]
				st.idx[node] = st.idx[node<<1|1]
			}
		}
	}
	build(1, 1, n)
	return st
}

func (st *segTree) apply(node int, v int64) {
	st.min[node] += v
	st.lazy[node] += v
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		st.apply(node<<1, st.lazy[node])
		st.apply(node<<1|1, st.lazy[node])
		st.lazy[node] = 0
	}
}

func (st *segTree) pull(node int) {
	if st.min[node<<1] <= st.min[node<<1|1] {
		st.min[node] = st.min[node<<1]
		st.idx[node] = st.idx[node<<1]
	} else {
		st.min[node] = st.min[node<<1|1]
		st.idx[node] = st.idx[node<<1|1]
	}
}

func (st *segTree) updateRange(node, l, r, L, R int, v int64) {
	if L > r || R < l {
		return
	}
	if L <= l && r <= R {
		st.apply(node, v)
		return
	}
	st.push(node)
	m := (l + r) >> 1
	st.updateRange(node<<1, l, m, L, R, v)
	st.updateRange(node<<1|1, m+1, r, L, R, v)
	st.pull(node)
}

func (st *segTree) updatePoint(node, l, r, p int) {
	if l == r {
		st.min[node] = INF
		st.lazy[node] = 0
		return
	}
	st.push(node)
	m := (l + r) >> 1
	if p <= m {
		st.updatePoint(node<<1, l, m, p)
	} else {
		st.updatePoint(node<<1|1, m+1, r, p)
	}
	st.pull(node)
}

func (st *segTree) queryMin() (int64, int) {
	return st.min[1], st.idx[1]
}

type inputD struct {
	n int
	v []int
	c []int64
	l []int64
	r []int64
}

func parseInputD(s string) (inputD, error) {
	rdr := bufio.NewReader(strings.NewReader(s))
	var n int
	if _, err := fmt.Fscan(rdr, &n); err != nil {
		return inputD{}, err
	}
	v := make([]int, n+1)
	c := make([]int64, n+1)
	l := make([]int64, n+1)
	r := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(rdr, &v[i], &c[i], &l[i], &r[i]); err != nil {
			return inputD{}, err
		}
	}
	return inputD{n: n, v: v, c: c, l: l, r: r}, nil
}

func solveD(inp inputD) []int {
	n := inp.n
	c := inp.c
	l := inp.l
	r := inp.r
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + c[i]
	}
	suff := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		suff[i] = suff[i+1] + c[i]
	}
	slack1 := make([]int64, n+1)
	slack2 := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		slack1[i] = pref[i-1] - l[i]
		slack2[i] = suff[i+1] - r[i]
	}
	st1 := newSegTree(slack1)
	st2 := newSegTree(slack2)
	alive := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		alive[i] = true
	}
	for {
		mn1, i1 := st1.queryMin()
		mn2, i2 := st2.queryMin()
		if mn1 >= 0 && mn2 >= 0 {
			break
		}
		var k int
		if mn1 < mn2 {
			k = i1
		} else {
			k = i2
		}
		if !alive[k] {
			if mn1 < 0 {
				st1.updatePoint(1, 1, n, i1)
			}
			if mn2 < 0 {
				st2.updatePoint(1, 1, n, i2)
			}
			continue
		}
		alive[k] = false
		if k+1 <= n {
			st1.updateRange(1, 1, n, k+1, n, -c[k])
		}
		if k-1 >= 1 {
			st2.updateRange(1, 1, n, 1, k-1, -c[k])
		}
		st1.updatePoint(1, 1, n, k)
		st2.updatePoint(1, 1, n, k)
	}
	var res []int
	for i := 1; i <= n; i++ {
		if alive[i] {
			res = append(res, i)
		}
	}
	return res
}

func verifyD(input, output string) error {
	inp, err := parseInputD(input)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}
	expected := solveD(inp)
	outFields := strings.Fields(strings.TrimSpace(output))
	if len(outFields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(outFields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k != len(expected) {
		return fmt.Errorf("expected %d trucks got %d", len(expected), k)
	}
	ids := make([]int, len(outFields)-1)
	for i := 1; i < len(outFields); i++ {
		val, err := strconv.Atoi(outFields[i])
		if err != nil {
			return fmt.Errorf("invalid id")
		}
		ids[i-1] = val
	}
	if len(ids) != k {
		return fmt.Errorf("expected %d ids got %d", k, len(ids))
	}
	for i := 0; i < k; i++ {
		if ids[i] != expected[i] {
			return fmt.Errorf("wrong id list")
		}
	}
	return nil
}

func runCase(bin, tc string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verifyD(tc, out.String())
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		v := rng.Intn(10) + 1
		c := rng.Intn(10) + 1
		l := rng.Intn(10)
		r := rng.Intn(10)
		fmt.Fprintf(&b, "%d %d %d %d\n", v, c, l, r)
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

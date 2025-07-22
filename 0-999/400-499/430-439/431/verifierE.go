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

func lowbit(x int) int { return x & -x }

func solveE(n int, h []int, queries []Query) string {
	disc := make([]int, len(h))
	copy(disc, h)
	for _, q := range queries {
		if q.typ == 1 {
			disc = append(disc, q.newv)
		}
	}
	sort.Ints(disc)
	m := 0
	for i := 0; i < len(disc); i++ {
		if i == 0 || disc[i] != disc[i-1] {
			disc[m] = disc[i]
			m++
		}
	}
	disc = disc[:m]
	cntFen := make([]int, m+1)
	sumFen := make([]int64, m+1)
	update := func(i, delta int) {
		d := int64(disc[i-1]) * int64(delta)
		for ; i <= m; i += lowbit(i) {
			cntFen[i] += delta
			sumFen[i] += d
		}
	}
	queryFen := func(i int) (cnt int, sum int64) {
		for ; i > 0; i -= lowbit(i) {
			cnt += cntFen[i]
			sum += sumFen[i]
		}
		return
	}
	for _, v := range h {
		id := sort.SearchInts(disc, v) + 1
		update(id, 1)
	}
	var out strings.Builder
	for _, qr := range queries {
		if qr.typ == 1 {
			idx := qr.idx - 1
			oldv := h[idx]
			tidOld := sort.SearchInts(disc, oldv) + 1
			update(tidOld, -1)
			tidNew := sort.SearchInts(disc, qr.newv) + 1
			update(tidNew, 1)
			h[idx] = qr.newv
		} else {
			v := qr.v
			l, r := 1, m
			nl := 1
			for l <= r {
				mid := (l + r) >> 1
				cnt, sum := queryFen(mid)
				if v+sum > int64(cnt)*int64(disc[mid-1]) {
					nl = mid
					l = mid + 1
				} else {
					r = mid - 1
				}
			}
			cnt, sum := queryFen(nl)
			avg := float64(v+sum) / float64(cnt)
			fmt.Fprintf(&out, "%f\n", avg)
		}
	}
	return out.String()
}

type Query struct {
	typ  int
	idx  int
	newv int
	v    int64
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	h := make([]int, n)
	for i := range h {
		h[i] = rng.Intn(10)
	}
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			idx := rng.Intn(n) + 1
			newv := rng.Intn(10)
			queries[i] = Query{typ: 1, idx: idx, newv: newv}
		} else {
			v := rng.Int63n(20)
			queries[i] = Query{typ: 2, v: v}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, qr := range queries {
		if qr.typ == 1 {
			fmt.Fprintf(&sb, "1 %d %d\n", qr.idx, qr.newv)
		} else {
			fmt.Fprintf(&sb, "2 %d\n", qr.v)
		}
	}
	return sb.String(), solveE(n, append([]int(nil), h...), queries)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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

	cases := []struct{ in, out string }{}
	for i := 0; i < 102; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.in, tc.out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

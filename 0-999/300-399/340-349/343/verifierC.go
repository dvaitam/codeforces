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

type testCase struct {
	n, m int
	h    []int64
	p    []int64
}

func canCover(h, p []int64, T int64) bool {
	n, m := len(h), len(p)
	j := 0
	for i := 0; i < n && j < m; i++ {
		hi := h[i]
		if p[j] > hi {
			limit := hi + T
			t := sort.Search(m, func(x int) bool { return p[x] > limit })
			j = t
		} else {
			l := p[j]
			if hi-l > T {
				return false
			}
			k := sort.Search(m, func(x int) bool { return p[x] > hi }) - 1
			tMax := k
			t1lim := T + 2*l - hi
			if t1lim >= l {
				idx := sort.Search(m, func(x int) bool { return p[x] > t1lim }) - 1
				if idx > tMax {
					tMax = idx
				}
			}
			rem := T - (hi - l)
			if rem >= 0 {
				t2lim := hi + rem/2
				idx2 := sort.Search(m, func(x int) bool { return p[x] > t2lim }) - 1
				if idx2 > tMax {
					tMax = idx2
				}
			}
			j = tMax + 1
		}
	}
	return j >= m
}

func expected(tc testCase) int64 {
	lo, hi := int64(0), int64(40000000000)
	for lo < hi {
		mid := (lo + hi) / 2
		if canCover(tc.h, tc.p, mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", tc.h[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", tc.p[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
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
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func generateCases(rng *rand.Rand) []testCase {
	cases := []testCase{}
	for len(cases) < 100 {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		hs := make([]int64, n)
		ps := make([]int64, m)
		vals := rand.Perm(50)
		for i := 0; i < n; i++ {
			hs[i] = int64(vals[i] + 1)
		}
		vals2 := rand.Perm(50)
		for i := 0; i < m; i++ {
			ps[i] = int64(vals2[i] + 1)
		}
		sort.Slice(hs, func(i, j int) bool { return hs[i] < hs[j] })
		sort.Slice(ps, func(i, j int) bool { return ps[i] < ps[j] })
		cases = append(cases, testCase{n, m, hs, ps})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

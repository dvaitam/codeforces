package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type pair struct {
	idx int
	cnt int64
}

type setData struct {
	l, r, total int64
	n           int
	a           []int64
	c           []int64
}

type testCaseC struct {
	m    int
	sets []setData
}

func expectedC(tc testCaseC) string {
	m := tc.m
	sets := tc.sets
	valMap := make(map[int64][]pair)
	uniq := make(map[int64]struct{})
	var L, R int64
	for i := 0; i < m; i++ {
		L += sets[i].l
		R += sets[i].r
		for j := 0; j < sets[i].n; j++ {
			v := sets[i].a[j]
			valMap[v] = append(valMap[v], pair{idx: i, cnt: sets[i].c[j]})
			uniq[v] = struct{}{}
		}
	}
	values := make([]int64, 0, len(uniq))
	for v := range uniq {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	prev := L - 1
	gap := false
	for _, v := range values {
		if v < L {
			continue
		}
		if v > R {
			break
		}
		if v > prev+1 {
			gap = true
			break
		}
		prev = v
	}
	if !gap && prev < R {
		gap = true
	}
	if gap {
		return "0"
	}
	baseCap := R - L
	best := int64(1<<63 - 1)
	for _, S := range values {
		if S < L || S > R {
			continue
		}
		F0 := int64(0)
		nonSCap := baseCap
		extraCap := int64(0)
		for _, p := range valMap[S] {
			set := sets[p.idx]
			t := set.total - p.cnt
			base := int64(0)
			if t < set.l {
				base = set.l - t
			}
			minrt := set.r
			if t < minrt {
				minrt = t
			}
			capi := int64(0)
			if minrt > set.l {
				capi = minrt - set.l
			}
			nonSCap += capi - (set.r - set.l)
			extra := set.r - set.l - capi
			remain := p.cnt - base
			if remain < extra {
				extra = remain
			}
			if extra < 0 {
				extra = 0
			}
			extraCap += extra
			F0 += base
		}
		need := int64(0)
		if S-L > nonSCap {
			need = S - L - nonSCap
		}
		if need <= extraCap {
			cand := F0 + need
			if cand < best {
				best = cand
			}
		}
	}
	return fmt.Sprint(best)
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 0, 100)
	for len(tests) < 100 {
		m := rand.Intn(3) + 1
		sets := make([]setData, m)
		for i := 0; i < m; i++ {
			n := rand.Intn(3) + 1
			a := make([]int64, n)
			used := map[int64]struct{}{}
			for j := 0; j < n; j++ {
				val := rand.Int63n(10) + 1
				for {
					if _, ok := used[val]; !ok {
						used[val] = struct{}{}
						break
					}
					val = rand.Int63n(10) + 1
				}
				a[j] = val
			}
			sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
			c := make([]int64, n)
			var total int64
			for j := 0; j < n; j++ {
				c[j] = int64(rand.Intn(3) + 1)
				total += c[j]
			}
			l := int64(rand.Intn(int(total)) + 1)
			r := l + int64(rand.Intn(int(total-l+1)))
			sets[i] = setData{l: l, r: r, total: total, n: n, a: a, c: c}
		}
		tests = append(tests, testCaseC{m: m, sets: sets})
	}
	return tests
}

func runCase(bin string, tc testCaseC) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.m))
	for _, s := range tc.sets {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", s.n, s.l, s.r))
		for i, v := range s.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range s.c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedC(tc)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

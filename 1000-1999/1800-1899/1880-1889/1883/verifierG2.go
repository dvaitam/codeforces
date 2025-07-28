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

type testCaseG2 struct {
	n int
	m int
	a []int
	b []int
}

func solveCaseG2(tc testCaseG2) string {
	n := tc.n
	a := append([]int(nil), tc.a...)
	b := append([]int(nil), tc.b...)
	sort.Ints(a)
	sort.Ints(b)
	pre := make([]int, n)
	for j := 0; j < n; j++ {
		pre[j] = sort.SearchInts(a, b[j])
	}
	hist := make([]int, n+1)
	used := 0
	for j := 0; j < n; j++ {
		if pre[j] > used {
			used++
		}
		hist[j+1] = used
	}
	k0 := hist[n]
	good := make([]int, n+1)
	exist := false
	for p := n; p >= 0; p-- {
		if p < n && pre[p] == hist[p] {
			exist = true
		}
		if exist {
			good[p] = 1
		}
	}
	counts := make([]int64, n+1)
	prev := 1
	for p := 0; p < n; p++ {
		nxt := tc.m
		if b[p]-1 < tc.m {
			nxt = b[p] - 1
		}
		if nxt >= prev {
			counts[p] = int64(nxt - prev + 1)
		}
		if b[p] > prev {
			prev = b[p]
		}
	}
	if prev <= tc.m {
		counts[n] = int64(tc.m - prev + 1)
	}
	var ans int64
	for p := 0; p <= n; p++ {
		if counts[p] > 0 {
			k := k0 + good[p]
			ops := n - k
			ans += counts[p] * int64(ops)
		}
	}
	return fmt.Sprint(ans)
}

func runCaseG2(bin string, tc testCaseG2) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCaseG2(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseG2(rng *rand.Rand) testCaseG2 {
	n := rng.Intn(5) + 2
	m := rng.Intn(10) + 1
	a := make([]int, n-1)
	for i := range a {
		a[i] = rng.Intn(m) + 1
	}
	b := make([]int, n)
	for i := range b {
		b[i] = rng.Intn(m) + 1
	}
	return testCaseG2{n: n, m: m, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseG2{{n: 2, m: 1, a: []int{1}, b: []int{1, 1}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseG2(rng))
	}
	for idx, tc := range cases {
		if err := runCaseG2(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

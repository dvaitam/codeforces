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

type testD struct {
	n int
	m int
	a []int64
	b []int64
}

func solveD(tc testD) int64 {
	a := append([]int64(nil), tc.a...)
	b := append([]int64(nil), tc.b...)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	prefA := make([]int64, len(a)+1)
	for i := range a {
		prefA[i+1] = prefA[i] + a[i]
	}
	prefB := make([]int64, len(b)+1)
	for i := range b {
		prefB[i+1] = prefB[i] + b[i]
	}
	totalB := prefB[len(b)]
	candidates := append(append([]int64{}, a...), b...)
	sort.Slice(candidates, func(i, j int) bool { return candidates[i] < candidates[j] })
	uniq := []int64{}
	var last int64 = -1 << 63
	for _, v := range candidates {
		if v != last {
			uniq = append(uniq, v)
			last = v
		}
	}
	ans := int64(-1)
	for _, v := range uniq {
		ia := sort.Search(len(a), func(i int) bool { return a[i] >= v })
		ib := sort.Search(len(b), func(i int) bool { return b[i] > v })
		costA := int64(v)*int64(ia) - prefA[ia]
		costB := totalB - prefB[ib] - int64(v)*int64(len(b)-ib)
		total := costA + costB
		if ans < 0 || total < ans {
			ans = total
		}
	}
	return ans
}

func genD(rng *rand.Rand) testD {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	a := make([]int64, n)
	b := make([]int64, m)
	for i := range a {
		a[i] = int64(rng.Intn(100) + 1)
	}
	for i := range b {
		b[i] = int64(rng.Intn(100) + 1)
	}
	return testD{n: n, m: m, a: a, b: b}
}

func runCase(bin string, tc testD) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testD{
		{n: 2, m: 2, a: []int64{1, 2}, b: []int64{3, 4}},
		{n: 1, m: 1, a: []int64{5}, b: []int64{1}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, genD(rng))
	}
	for i, tc := range cases {
		expect := solveD(tc)
		out, err := runCase(exe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if fmt.Sprint(expect) != out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

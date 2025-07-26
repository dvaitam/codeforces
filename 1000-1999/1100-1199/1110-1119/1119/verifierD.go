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
	a  []int64
	lr [][2]int64
}

func expected(tc testCase) string {
	n := len(tc.a)
	a := append([]int64(nil), tc.a...)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	if n <= 1 {
		var res strings.Builder
		for i := range tc.lr {
			if i > 0 {
				res.WriteByte(' ')
			}
			res.WriteString("0")
		}
		return res.String()
	}
	b := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		b[i] = a[i+1] - a[i]
	}
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	ps := make([]int64, n)
	for i := 1; i < n; i++ {
		ps[i] = ps[i-1] + b[i-1]
	}
	var out strings.Builder
	for qi, lr := range tc.lr {
		if qi > 0 {
			out.WriteByte(' ')
		}
		k := lr[1] - lr[0] + 1
		idx := sort.Search(len(b), func(i int) bool { return b[i] > k })
		ans := ps[idx] + int64(n-idx)*k
		out.WriteString(fmt.Sprintf("%d", ans))
	}
	return out.String()
}

func buildCase(a []int64, lr [][2]int64) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(lr)))
	for _, p := range lr {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return sb.String(), expected(testCase{a: a, lr: lr})
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(50))
	}
	q := rng.Intn(10) + 1
	lr := make([][2]int64, q)
	for i := 0; i < q; i++ {
		l := int64(rng.Intn(50))
		r := l + int64(rng.Intn(50))
		lr[i] = [2]int64{l, r}
	}
	return buildCase(a, lr)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	in, exp := buildCase([]int64{1, 2, 4}, [][2]int64{{0, 1}, {1, 3}})
	cases = append(cases, in)
	exps = append(exps, exp)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 102 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}
	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

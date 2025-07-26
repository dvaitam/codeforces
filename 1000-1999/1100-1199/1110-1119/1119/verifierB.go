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
	h int64
	a []int64
}

func expected(tc testCase) string {
	n := len(tc.a)
	h := tc.h
	result := 0
	for i := 1; i <= n; i++ {
		b := make([]int64, i)
		copy(b, tc.a[:i])
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		var sum int64
		for j := 0; j < i; j += 2 {
			sum += b[j]
		}
		if sum <= h {
			result = i
		} else {
			break
		}
	}
	return fmt.Sprintf("%d", result)
}

func buildCase(a []int64, h int64) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(a), h))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(testCase{h: h, a: a})
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	h := int64(rng.Intn(100) + 1)
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(int(h)) + 1)
	}
	return buildCase(a, h)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	in, exp := buildCase([]int64{1, 2, 3}, 3)
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

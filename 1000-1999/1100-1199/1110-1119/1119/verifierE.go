package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	a []int64
}

func expected(tc testCase) string {
	var leftover int64
	var ans int64
	for _, x := range tc.a {
		pairs := leftover
		if x/2 < pairs {
			pairs = x / 2
		}
		rem := x - pairs*2
		triples := rem / 3
		ans += pairs + triples
		leftover = leftover + x - 3*(pairs+triples)
	}
	return fmt.Sprintf("%d", ans)
}

func buildCase(a []int64) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(testCase{a: a})
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(100))
	}
	return buildCase(a)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	in, exp := buildCase([]int64{1, 2, 3})
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

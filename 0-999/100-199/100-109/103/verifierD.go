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

type query struct{ a, b int }

type testCase struct {
	input  string
	expect []int64
}

func solve(w []int64, qs []query) []int64 {
	res := make([]int64, len(qs))
	for i, q := range qs {
		var sum int64
		for j := q.a; j < len(w); j += q.b {
			sum += w[j]
		}
		res[i] = sum
	}
	return res
}

func buildCase(w []int64, qs []query) testCase {
	var sb strings.Builder
	n := len(w)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", w[i])
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(qs))
	for _, q := range qs {
		fmt.Fprintf(&sb, "%d %d\n", q.a+1, q.b)
	}
	return testCase{input: sb.String(), expect: solve(w, qs)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	w := make([]int64, n)
	for i := range w {
		w[i] = rng.Int63n(100) - 50
	}
	p := rng.Intn(20) + 1
	qs := make([]query, p)
	for i := range qs {
		qs[i].a = rng.Intn(n)
		qs[i].b = rng.Intn(n) + 1
	}
	return buildCase(w, qs)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.expect) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.expect), len(fields))
	}
	for i, f := range fields {
		var v int64
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != tc.expect[i] {
			return fmt.Errorf("expected %v got %v", tc.expect, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int64{1}, []query{{0, 1}}))
	cases = append(cases, buildCase([]int64{5, 3, 2}, []query{{0, 1}, {1, 2}}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

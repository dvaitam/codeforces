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

type query struct {
	l, r int
}

type testCase struct {
	input    string
	expected []int64
}

func solveE(p []int, queries []query) []int64 {
	res := make([]int64, len(queries))
	for idx, q := range queries {
		var count int64
		for i := q.l; i <= q.r; i++ {
			minVal := p[i]
			maxVal := p[i]
			for j := i; j <= q.r; j++ {
				if p[j] < minVal {
					minVal = p[j]
				}
				if p[j] > maxVal {
					maxVal = p[j]
				}
				if maxVal-minVal == j-i {
					count++
				}
			}
		}
		res[idx] = count
	}
	return res
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	p := make([]int, n)
	for i := range p {
		p[i] = rng.Intn(5) + 1
	}
	qcnt := rng.Intn(5) + 1
	queries := make([]query, qcnt)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", p[i])
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", qcnt)
	for i := 0; i < qcnt; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		queries[i] = query{l: l - 1, r: r - 1}
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	expected := solveE(p, queries)
	return testCase{input: sb.String(), expected: expected}
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
	if len(fields) != len(tc.expected) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.expected), len(fields))
	}
	for i, f := range fields {
		var v int64
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != tc.expected[i] {
			return fmt.Errorf("expected %v got %v", tc.expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

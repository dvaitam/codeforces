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
	input    string
	expected []int64
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

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(20) + 1
	base := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		base[i] = int64(rng.Intn(1000) + 1)
	}
	initVals := append([]int64(nil), base...)
	add := int64(0)
	var expect []int64
	ops := make([]string, 0, m+1)
	for i := 0; i < m; i++ {
		t := rng.Intn(3) + 1
		switch t {
		case 1:
			idx := rng.Intn(n) + 1
			x := int64(rng.Intn(1000) + 1)
			ops = append(ops, fmt.Sprintf("1 %d %d", idx, x))
			base[idx] = x - add
		case 2:
			y := int64(rng.Intn(100) + 1)
			ops = append(ops, fmt.Sprintf("2 %d", y))
			add += y
		case 3:
			q := rng.Intn(n) + 1
			ops = append(ops, fmt.Sprintf("3 %d", q))
			expect = append(expect, base[q]+add)
		}
	}
	if len(expect) == 0 {
		q := 1
		ops = append(ops, fmt.Sprintf("3 %d", q))
		expect = append(expect, base[q]+add)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", initVals[i]))
	}
	sb.WriteString("\n")
	for _, op := range ops {
		sb.WriteString(op)
		sb.WriteString("\n")
	}
	return testCase{input: sb.String(), expected: expect}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{
			input:    "1 1\n5\n3 1\n",
			expected: []int64{5},
		},
		{
			input:    "2 3\n1 2\n2 1\n3 2\n2 3\n",
			expected: []int64{3},
		},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

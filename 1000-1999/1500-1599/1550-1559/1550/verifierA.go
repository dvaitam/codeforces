package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	sVals    []int
	expected []int
}

func ceilSqrt(x int) int {
	if x <= 0 {
		return 0
	}
	r := int(math.Sqrt(float64(x)))
	if r*r < x {
		r++
	}
	return r
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(10) + 1
	sVals := make([]int, t)
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		val := rng.Intn(5000) + 1
		sVals[i] = val
		expected[i] = ceilSqrt(val)
	}
	return testCase{sVals: sVals, expected: expected}
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.sVals)))
	for _, v := range tc.sVals {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func (tc testCase) output() string {
	var sb strings.Builder
	for i, v := range tc.expected {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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
	n, m   int
	minVal int
	maxVal int
	temps  []int
}

func expected(tc testCase) string {
	hasMin := false
	hasMax := false
	for _, t := range tc.temps {
		if t < tc.minVal || t > tc.maxVal {
			return "Incorrect"
		}
		if t == tc.minVal {
			hasMin = true
		}
		if t == tc.maxVal {
			hasMax = true
		}
	}
	need := 0
	if !hasMin {
		need++
	}
	if !hasMax {
		need++
	}
	if need <= tc.n-tc.m {
		return "Correct"
	}
	return "Incorrect"
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.minVal, tc.maxVal))
	for i := 0; i < tc.m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(tc.temps[i]))
	}
	if tc.m > 0 {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	got := strings.TrimSpace(out.String())
	exp := expected(tc)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	m := rng.Intn(n-1) + 1
	minVal := rng.Intn(98) + 1
	maxVal := rng.Intn(100-minVal) + minVal + 1
	temps := make([]int, m)
	for i := 0; i < m; i++ {
		temps[i] = rng.Intn(100) + 1
	}
	return testCase{n: n, m: m, minVal: minVal, maxVal: maxVal, temps: temps}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, testCase{n: 2, m: 1, minVal: 1, maxVal: 2, temps: []int{1}})
	tests = append(tests, testCase{n: 3, m: 2, minVal: 1, maxVal: 3, temps: []int{1, 3}})
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validateOutput(n int, a, b []int64, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %q", fields[0])
	}
	limit := n/2 + 1
	if k < 1 || k > limit {
		return fmt.Errorf("k out of range: got %d, allowed [1,%d]", k, limit)
	}
	if len(fields) != 1+k {
		return fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}

	seen := make(map[int]bool, k)
	var sumSelA, sumSelB int64
	var sumAllA, sumAllB int64
	for i := 1; i <= n; i++ {
		sumAllA += a[i]
		sumAllB += b[i]
	}

	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return fmt.Errorf("invalid index token: %q", fields[1+i])
		}
		if idx < 1 || idx > n {
			return fmt.Errorf("index out of range: %d", idx)
		}
		if seen[idx] {
			return fmt.Errorf("duplicate index: %d", idx)
		}
		seen[idx] = true
		sumSelA += a[idx]
		sumSelB += b[idx]
	}

	if 2*sumSelA <= sumAllA {
		return fmt.Errorf("A inequality failed: 2*%d <= %d", sumSelA, sumAllA)
	}
	if 2*sumSelB <= sumAllB {
		return fmt.Errorf("B inequality failed: 2*%d <= %d", sumSelB, sumAllB)
	}
	return nil
}

type testCase struct {
	a []int64
	b []int64
}

func buildInput(tc testCase) string {
	n := len(tc.a)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.a[i], 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.b[i], 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func edgeCases() []testCase {
	return []testCase{
		{a: []int64{1}, b: []int64{1}},
		{a: []int64{1, 2}, b: []int64{2, 1}},
		{a: []int64{5, 1, 1}, b: []int64{1, 5, 1}},
		{a: []int64{1, 1, 100, 100}, b: []int64{100, 100, 1, 1}},
		{a: []int64{1, 2, 3, 4, 5}, b: []int64{5, 4, 3, 2, 1}},
		{a: []int64{1_000_000_000, 1, 1, 1, 1}, b: []int64{1, 1_000_000_000, 1, 1, 1}},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(1_000_000_000) + 1)
		b[i] = int64(rng.Intn(1_000_000_000) + 1)
	}
	return testCase{a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	tests := edgeCases()
	for i := 0; i < 120; i++ {
		tests = append(tests, randomCase(rng))
	}

	for idx, tc := range tests {
		testNum := idx + 1
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", testNum, err)
			os.Exit(1)
		}
		n := len(tc.a)
		aVals := make([]int64, n+1)
		bVals := make([]int64, n+1)
		for i := 0; i < n; i++ {
			aVals[i+1] = tc.a[i]
			bVals[i+1] = tc.b[i]
		}
		if err := validateOutput(n, aVals, bVals, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s got: %s\n", testNum, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

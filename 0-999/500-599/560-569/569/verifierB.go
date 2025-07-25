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
	n   int
	arr []int
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func minimalChanges(tc testCase) int {
	seen := make(map[int]struct{})
	for _, v := range tc.arr {
		if v >= 1 && v <= tc.n {
			seen[v] = struct{}{}
		}
	}
	return tc.n - len(seen)
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	res := make([]int, n)
	for i, f := range fields {
		var x int
		if _, err := fmt.Sscan(f, &x); err != nil {
			return nil, fmt.Errorf("bad integer %q", f)
		}
		res[i] = x
	}
	return res, nil
}

func isPermutation(arr []int, n int) bool {
	seen := make([]bool, n+1)
	for _, v := range arr {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	ans, err := parseOutput(out, tc.n)
	if err != nil {
		return err
	}
	if !isPermutation(ans, tc.n) {
		return fmt.Errorf("output is not a permutation")
	}
	changes := 0
	for i := 0; i < tc.n; i++ {
		if ans[i] != tc.arr[i] {
			changes++
		}
	}
	if changes != minimalChanges(tc) {
		return fmt.Errorf("expected %d changes, got %d", minimalChanges(tc), changes)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Float64() < 0.5 {
			arr[i] = rng.Intn(n) + 1
		} else {
			arr[i] = rng.Intn(1000) + n + 1
		}
	}
	return testCase{n: n, arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, testCase{n: 3, arr: []int{1, 2, 3}})
	cases = append(cases, testCase{n: 3, arr: []int{1, 1, 1}})
	cases = append(cases, testCase{n: 5, arr: []int{1, 2, 2, 3, 3}})
	cases = append(cases, testCase{n: 3, arr: []int{4, 5, 6}})
	cases = append(cases, testCase{n: 1, arr: []int{1}})
	for len(cases) < 105 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n, m int
}

func runCase(bin string, tc testCase) ([]int, error) {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != tc.m {
		return nil, fmt.Errorf("expected %d numbers, got %d", tc.m, len(fields))
	}
	res := make([]int, tc.m)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func verifyCase(tc testCase, arr []int) error {
	if len(arr) != tc.m {
		return fmt.Errorf("expected %d numbers, got %d", tc.m, len(arr))
	}
	sum := 0
	minv := arr[0]
	maxv := arr[0]
	for _, v := range arr {
		if v <= 0 {
			return fmt.Errorf("non-positive value %d", v)
		}
		sum += v
		if v < minv {
			minv = v
		}
		if v > maxv {
			maxv = v
		}
	}
	if sum != tc.n {
		return fmt.Errorf("sum mismatch: expected %d got %d", tc.n, sum)
	}
	diff := maxv - minv
	expectDiff := 0
	if tc.n%tc.m != 0 {
		expectDiff = 1
	}
	if diff != expectDiff {
		return fmt.Errorf("expected max-min difference %d, got %d", expectDiff, diff)
	}
	return nil
}

func generateCases() []testCase {
	cases := []testCase{{1, 1}, {2, 1}, {2, 2}, {3, 2}, {100, 1}, {100, 100}, {7, 3}, {10, 10}, {99, 50}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		n := rng.Intn(100) + 1
		m := rng.Intn(n) + 1
		cases = append(cases, testCase{n, m})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		arr, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d)\n", i+1, err, tc.n, tc.m)
			os.Exit(1)
		}
		if err := verifyCase(tc, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d)\n", i+1, err, tc.n, tc.m)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func solveReference(input string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &a[i]); err != nil {
			return 0, err
		}
	}
	if n == 0 {
		return 0, nil
	}
	candidates := make(map[int64]struct{})
	if a[0] > 0 {
		candidates[a[0]-1] = struct{}{}
	}
	limit := int64(100000)
	if a[0]-1 <= limit {
		for x := int64(0); x < a[0]; x++ {
			candidates[x] = struct{}{}
		}
	} else {
		for _, v := range a[1:] {
			if v == 0 {
				continue
			}
			for k := int64(1); ; k++ {
				val := k*v - 1
				if val >= a[0] || val < 0 {
					break
				}
				candidates[val] = struct{}{}
				if k*v-1 > limit {
					break
				}
			}
		}
		for i := int64(0); i < limit && a[0]-1-i >= 0; i++ {
			candidates[a[0]-1-i] = struct{}{}
		}
	}
	var best int64
	for x := range candidates {
		if x < 0 || x >= a[0] {
			continue
		}
		cur := x
		var sum int64
		for i := 0; i < n; i++ {
			if a[i] == 0 {
				cur = 0
			} else {
				cur = cur % a[i]
			}
			sum += cur
		}
		if sum > best {
			best = sum
		}
	}
	return best, nil
}

func makeCase(name string, arr []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_small", []int64{1}),
		makeCase("single_large", []int64{999999999999}),
		makeCase("increasing_small", []int64{3, 5, 7}),
		makeCase("mixed_medium", []int64{10, 4, 8, 6}),
		makeCase("descending", []int64{10, 9, 8, 7, 6}),
		makeCase("two_values", []int64{5, 2}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 120; i++ {
		n := rng.Intn(6) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(100) + 1)
		}
		tests = append(tests, makeCase(fmt.Sprintf("small_rand_%d", i+1), arr))
	}
	for i := 0; i < 40; i++ {
		n := rng.Intn(50) + 10
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(1000000) + 1)
		}
		tests = append(tests, makeCase(fmt.Sprintf("medium_rand_%d", i+1), arr))
	}
	for i := 0; i < 20; i++ {
		n := rng.Intn(100) + 50
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Int63n(1_000_000_000_000) + 1
		}
		tests = append(tests, makeCase(fmt.Sprintf("large_rand_%d", i+1), arr))
	}
	return tests
}

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (int64, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %v", fields)
	}
	var val int64
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveReference(tc.input)
		if err != nil {
			fmt.Printf("failed to solve reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) failed: expect %d, got %d\ninput:\n%s\n", idx+1, tc.name, expect, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

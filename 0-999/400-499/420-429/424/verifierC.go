package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	name  string
	input string
}

func xor0to(x int64) int64 {
	switch x & 3 {
	case 0:
		return x
	case 1:
		return 1
	case 2:
		return x + 1
	default:
		return 0
	}
}

func solveRef(input string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	p := make([]int64, n)
	for i := int64(0); i < n; i++ {
		if _, err := fmt.Fscan(reader, &p[i]); err != nil {
			return 0, err
		}
	}
	var pxor int64
	for i := int64(0); i+1 < n; i++ {
		pxor ^= p[i]
	}
	var t int64
	m := n - 1
	for k := int64(2); k <= n; k++ {
		full := m / k
		rem := m % k
		if full&1 == 1 {
			t ^= xor0to(k - 1)
		}
		t ^= xor0to(rem)
	}
	return pxor ^ t, nil
}

func makeCase(name string, arr []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_element_zero", []int64{0}),
		makeCase("single_element_large", []int64{123456789}),
		makeCase("two_elements", []int64{5, 7}),
		makeCase("three_elements", []int64{1, 2, 3}),
		makeCase("increasing_large", []int64{1, 3, 5, 7, 9, 11}),
		makeCase("random_small", []int64{8, 6, 7, 5, 3, 0, 9}),
		makeCase("alternating_bits", []int64{1, 2, 4, 8, 16, 32, 64, 128}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(424))
	var tests []testCase
	gen := func(prefix string, cnt, maxN int, maxVal int64) {
		for i := 0; i < cnt; i++ {
			n := rng.Intn(maxN) + 1
			arr := make([]int64, n)
			for j := 0; j < n; j++ {
				arr[j] = rng.Int63n(maxVal + 1)
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), arr))
		}
	}
	gen("small", 120, 10, 100)
	gen("medium", 120, 200, 100000)
	return tests
}

func largeTests() []testCase {
	rng := rand.New(rand.NewSource(424424))
	var tests []testCase
	for i := 0; i < 5; i++ {
		n := 100000 + rng.Intn(200000)
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = rng.Int63n(2_000_000_000)
		}
		tests = append(tests, makeCase(fmt.Sprintf("large_%d", i+1), arr))
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseSingleInteger(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no output produced")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("expected single integer, got extra data: %v", fields)
	}
	var val int64
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(append(handcraftedTests(), randomTests()...), largeTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		gotStr, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		gotVal, parseErr := parseSingleInteger(gotStr)
		if parseErr != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, parseErr, tc.input, gotStr)
			os.Exit(1)
		}
		if gotVal != expect {
			fmt.Printf("test %d (%s) failed: expect %d, got %d\ninput:\n%s\n", idx+1, tc.name, expect, gotVal, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

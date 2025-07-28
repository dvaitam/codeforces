package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input    string
	expected int
}

func solveCase(n int, arr []int) int {
	zeros := 0
	freq := make(map[int]int)
	dup := false
	for _, v := range arr {
		if v == 0 {
			zeros++
		}
		freq[v]++
		if freq[v] > 1 {
			dup = true
		}
	}
	if zeros > 0 {
		return n - zeros
	}
	if dup {
		return n
	}
	return n + 1
}

func buildCase(arr []int) testCase {
	var sb strings.Builder
	n := len(arr)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solveCase(n, arr)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(99) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(101)
	}
	return buildCase(arr)
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
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	var val int
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	var cases []testCase

	// deterministic edge cases
	cases = append(cases, buildCase([]int{0, 0}))
	cases = append(cases, buildCase([]int{1, 2}))
	cases = append(cases, buildCase([]int{1, 1}))

	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

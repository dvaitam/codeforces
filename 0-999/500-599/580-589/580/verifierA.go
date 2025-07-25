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
	expected string
}

func solve(nums []int) int {
	maxLen, currLen := 0, 0
	prev := 0
	for i, x := range nums {
		if i == 0 || x >= prev {
			currLen++
		} else {
			currLen = 1
		}
		if currLen > maxLen {
			maxLen = currLen
		}
		prev = x
	}
	return maxLen
}

func buildCase(nums []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(nums)))
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	ans := solve(nums)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d\n", ans)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(1000)
	}
	return buildCase(nums)
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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

	cases := []testCase{
		buildCase([]int{1}),
		buildCase([]int{1, 2, 3}),
		buildCase([]int{3, 2, 1}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

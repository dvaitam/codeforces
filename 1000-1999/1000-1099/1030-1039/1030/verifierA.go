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
	input    string
	expected string
}

func solve(nums []int) string {
	for _, v := range nums {
		if v == 1 {
			return "HARD\n"
		}
	}
	return "EASY\n"
}

func buildCase(nums []int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(nums)))
	sb.WriteByte('\n')
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solve(nums)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(2)
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
	got := out.String()
	if got != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, got)
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
		buildCase([]int{0}),
		buildCase([]int{1, 0, 0}),
	}
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

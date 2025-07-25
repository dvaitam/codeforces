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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxLen := 1
	cur := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[i-1] {
			cur++
		} else {
			if cur > maxLen {
				maxLen = cur
			}
			cur = 1
		}
	}
	if cur > maxLen {
		maxLen = cur
	}
	return maxLen
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
	expect := strconv.Itoa(solve(nums))
	return testCase{input: sb.String(), expected: expect}
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(100)
	}
	return buildCase(nums)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int{1}))
	cases = append(cases, buildCase([]int{1, 2, 3}))
	cases = append(cases, buildCase([]int{3, 2, 1}))
	cases = append(cases, buildCase([]int{1, 3, 2, 4, 5}))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

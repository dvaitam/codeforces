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

func solve(nums []int64) int64 {
	counts := make(map[int64]int)
	var result int64
	for _, a := range nums {
		for p := 0; p < 32; p++ {
			need := (int64(1) << p) - a
			if c, ok := counts[need]; ok {
				result += int64(c)
			}
		}
		counts[a]++
	}
	return result
}

func buildCase(nums []int64) testCase {
	var sb strings.Builder
	n := len(nums)
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	expect := strconv.FormatInt(solve(nums), 10)
	return testCase{input: sb.String(), expected: expect}
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	nums := make([]int64, n)
	for i := range nums {
		nums[i] = int64(rng.Intn(100))
	}
	return buildCase(nums)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int64{1, 3, 5, 7}))
	cases = append(cases, buildCase([]int64{1, 1, 1, 1}))
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

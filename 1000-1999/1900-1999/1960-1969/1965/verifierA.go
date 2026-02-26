package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(nums []int) string {
	seen := make(map[int]struct{})
	var unique []int
	for _, v := range nums {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			unique = append(unique, v)
		}
	}
	sort.Ints(unique)
	for i, v := range unique {
		if v > i+1 {
			if i%2 == 0 {
				return "Alice"
			}
			return "Bob"
		}
	}
	if len(unique)%2 == 1 {
		return "Alice"
	}
	return "Bob"
}

func buildCase(nums []int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(nums)))
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solve(nums)}
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(10) + 1
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
	cases = append(cases, buildCase([]int{1, 1}))
	cases = append(cases, buildCase([]int{1, 2}))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}
	for idx, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

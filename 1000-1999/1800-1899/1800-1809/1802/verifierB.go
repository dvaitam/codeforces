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
	plan   []int
	input  string
	expect int
}

func maxCagesKnown(n int) int {
	if n == 0 {
		return 0
	}
	return (n + 2) / 2
}

func solveCase(plan []int) int {
	known := 0
	unknown := 0
	ans := 0
	current := func() int { return maxCagesKnown(known) + unknown }
	ans = current()
	for _, b := range plan {
		if b == 1 {
			unknown++
		} else {
			known += unknown
			unknown = 0
		}
		if c := current(); c > ans {
			ans = c
		}
	}
	return ans
}

func buildCase(plan []int) testCase {
	n := len(plan)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range plan {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{plan: plan, input: sb.String(), expect: solveCase(plan)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	plan := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			plan[i] = 1
		} else {
			plan[i] = 2
		}
	}
	return buildCase(plan)
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

func checkOutput(out string, tc testCase) error {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 1 {
		return fmt.Errorf("expected single number")
	}
	var val int
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase([]int{1}))
	cases = append(cases, buildCase([]int{2}))
	cases = append(cases, buildCase([]int{1, 1, 2, 2}))
	cases = append(cases, buildCase([]int{1, 2, 1, 2}))
	cases = append(cases, buildCase([]int{2, 2, 2}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d mismatch: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

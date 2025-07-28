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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(nums []int) string {
	var sb strings.Builder
	for _, n := range nums {
		if n <= 4 {
			sb.WriteString("Bob\n")
		} else {
			sb.WriteString("Alice\n")
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1
	var input strings.Builder
	var expected strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", t))
	nums := make([]int, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(99) + 2
		nums[i] = n
		input.WriteString(fmt.Sprintf("%d\n", n))
	}
	expected.WriteString(solveA(nums))
	return testCase{input: input.String(), expected: expected.String()}
}

func runCase(bin string, tc testCase) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
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
	cases := []testCase{{input: "1\n5\n", expected: "Alice\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

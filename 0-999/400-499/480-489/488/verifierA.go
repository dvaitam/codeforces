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

func has8(x int64) bool {
	if x < 0 {
		x = -x
	}
	for x > 0 {
		if x%10 == 8 {
			return true
		}
		x /= 10
	}
	return false
}

func solve(a int64) string {
	for i := int64(1); ; i++ {
		if has8(a + i) {
			return fmt.Sprintf("%d", i)
		}
	}
}

func generateCase(rng *rand.Rand) testCase {
	a := rng.Int63n(2000000001) - 1000000000
	input := fmt.Sprintf("%d\n", a)
	expected := solve(a)
	return testCase{input: input, expected: expected}
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
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	var cases []testCase
	// simple deterministic case
	cases = append(cases, testCase{input: "7\n", expected: solve(7)})
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

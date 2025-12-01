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
	g int
	c int
	l int
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{88, 94, 95},
		{100, 80, 81},
		{98, 99, 98},
		{95, 86, 85},
		{80, 80, 80},
		{100, 100, 100},
	}
}

func randomTest(rng *rand.Rand) testCase {
	return testCase{
		g: rng.Intn(21) + 80,
		c: rng.Intn(21) + 80,
		l: rng.Intn(21) + 80,
	}
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n", tc.g, tc.c, tc.l)
}

func expected(tc testCase) string {
	values := []int{tc.g, tc.c, tc.l}
	minVal, maxVal := values[0], values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal-minVal >= 10 {
		return "check again"
	}
	// median
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
	return fmt.Sprintf("final %d", values[1])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		out, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput: %s", idx+1, err, input)
			os.Exit(1)
		}
		exp := expected(tc)
		if out != exp {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %q got %q\ninput: %s", idx+1, exp, out, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

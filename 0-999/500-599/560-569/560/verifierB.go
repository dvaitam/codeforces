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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canFit(A, B, a2, b2, a3, b3 int) bool {
	for r1 := 0; r1 < 2; r1++ {
		x1, y1 := a2, b2
		if r1 == 1 {
			x1, y1 = b2, a2
		}
		for r2 := 0; r2 < 2; r2++ {
			x2, y2 := a3, b3
			if r2 == 1 {
				x2, y2 = b3, a3
			}
			if x1+x2 <= A && max(y1, y2) <= B {
				return true
			}
			if y1+y2 <= B && max(x1, x2) <= A {
				return true
			}
		}
	}
	return false
}

func buildCase(a1, b1, a2, b2, a3, b3 int) testCase {
	input := fmt.Sprintf("%d %d\n%d %d\n%d %d\n", a1, b1, a2, b2, a3, b3)
	expect := "NO"
	if canFit(a1, b1, a2, b2, a3, b3) || canFit(b1, a1, a2, b2, a3, b3) {
		expect = "YES"
	}
	return testCase{input: input, expected: expect}
}

func generateRandomCase(rng *rand.Rand) testCase {
	a1 := rng.Intn(1000) + 1
	b1 := rng.Intn(1000) + 1
	a2 := rng.Intn(1000) + 1
	b2 := rng.Intn(1000) + 1
	a3 := rng.Intn(1000) + 1
	b3 := rng.Intn(1000) + 1
	return buildCase(a1, b1, a2, b2, a3, b3)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// deterministic cases
	cases = append(cases, buildCase(2, 2, 1, 2, 1, 1)) // YES
	cases = append(cases, buildCase(3, 3, 2, 2, 2, 2)) // NO
	cases = append(cases, buildCase(5, 4, 3, 4, 2, 2)) // YES
	cases = append(cases, buildCase(1, 1, 1, 1, 1, 1)) // NO

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

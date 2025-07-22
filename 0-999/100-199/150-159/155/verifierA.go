package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type TestCase struct {
	input  string
	output string
}

func computeAmazing(scores []int) int {
	if len(scores) <= 1 {
		return 0
	}
	best, worst := scores[0], scores[0]
	count := 0
	for i := 1; i < len(scores); i++ {
		if scores[i] > best {
			count++
			best = scores[i]
		} else if scores[i] < worst {
			count++
			worst = scores[i]
		}
	}
	return count
}

func generateTests() []TestCase {
	rnd := rand.New(rand.NewSource(42))
	tests := make([]TestCase, 0, 110)
	// Deterministic edge cases
	tests = append(tests, TestCase{input: "1\n0\n", output: "0"})
	tests = append(tests, TestCase{input: "5\n1 2 3 4 5\n", output: "4"})
	tests = append(tests, TestCase{input: "5\n5 4 3 2 1\n", output: "4"})
	tests = append(tests, TestCase{input: "3\n10 10 10\n", output: "0"})

	for i := 0; i < 100; i++ {
		n := rnd.Intn(20) + 1 // 1..20
		scores := make([]int, n)
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			scores[j] = rnd.Intn(1000)
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%d", scores[j]))
		}
		b.WriteByte('\n')
		out := computeAmazing(scores)
		tests = append(tests, TestCase{input: b.String(), output: fmt.Sprintf("%d", out)})
	}
	sort.Slice(tests, func(i, j int) bool { return i < j })
	return tests
}

func runTest(binary string, tc TestCase) error {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %v", err)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.output)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runTest(binary, tc); err != nil {
			fmt.Printf("Test %d failed: %v\nInput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}

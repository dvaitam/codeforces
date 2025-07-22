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

type Card struct {
	a int
	b int
}

func computeMaxPoints(cards []Card) int64 {
	var zero []int
	var sumA1 int64
	var extra int64
	for _, c := range cards {
		if c.b == 0 {
			zero = append(zero, c.a)
		} else {
			sumA1 += int64(c.a)
			extra += int64(c.b - 1)
		}
	}
	sort.Slice(zero, func(i, j int) bool { return zero[i] > zero[j] })
	k0 := int(extra + 1)
	if k0 > len(zero) {
		k0 = len(zero)
	}
	var sumA0 int64
	for i := 0; i < k0; i++ {
		sumA0 += int64(zero[i])
	}
	return sumA1 + sumA0
}

func generateTests() []TestCase {
	rnd := rand.New(rand.NewSource(99))
	tests := make([]TestCase, 0, 110)
	// Edge cases
	tests = append(tests, TestCase{input: "1\n0 0\n", output: "0"})
	tests = append(tests, TestCase{input: "1\n5 1\n", output: "5"})
	tests = append(tests, TestCase{input: "2\n1 0\n2 0\n", output: "2"})
	tests = append(tests, TestCase{input: "2\n1 1\n2 0\n", output: "3"})
	for i := 0; i < 100; i++ {
		n := rnd.Intn(10) + 1 // 1..10
		cards := make([]Card, n)
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			a := rnd.Intn(20)
			bVal := rnd.Intn(4) // 0..3
			cards[j] = Card{a: a, b: bVal}
			b.WriteString(fmt.Sprintf("%d %d\n", a, bVal))
		}
		res := computeMaxPoints(cards)
		tests = append(tests, TestCase{input: b.String(), output: fmt.Sprintf("%d", res)})
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2171G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2171G.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(target string) *exec.Cmd {
	switch filepath.Ext(target) {
	case ".go":
		return exec.Command("go", "run", target)
	case ".py":
		return exec.Command("python3", target)
	default:
		return exec.Command(target)
	}
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func sampleTests() []string {
	return []string{
		"88 94 95\n",
		"100 80 81\n",
		"98 99 98\n",
		"95 86 85\n",
	}
}

func deterministicTests() []string {
	var tests []string
	// Exhaust all combinations where each score is one of a few representative values.
	values := []int{80, 85, 90, 95, 100}
	for _, g := range values {
		for _, c := range values {
			for _, l := range values {
				tests = append(tests, fmt.Sprintf("%d %d %d\n", g, c, l))
			}
		}
	}
	// Exhaust all triples where scores differ by exactly 10.
	for base := 80; base <= 90; base++ {
		tests = append(tests, fmt.Sprintf("%d %d %d\n", base, base+10, base+5))
		tests = append(tests, fmt.Sprintf("%d %d %d\n", base, base+10, base))
	}
	return tests
}

func randomTests(count int) []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, count)
	for i := 0; i < count; i++ {
		g := rng.Intn(21) + 80
		c := rng.Intn(21) + 80
		l := rng.Intn(21) + 80
		tests[i] = fmt.Sprintf("%d %d %d\n", g, c, l)
	}
	return tests
}

func buildTests() []string {
	var tests []string
	tests = append(tests, sampleTests()...)
	tests = append(tests, deterministicTests()...)
	tests = append(tests, randomTests(200)...)
	return tests
}

func normalizeOutput(out string) string {
	return strings.TrimSpace(out)
}

func compareOutputs(expected, got string) error {
	if expected != got {
		return fmt.Errorf("expected %q, got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, input := range tests {
		expectRaw, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotRaw, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expect := normalizeOutput(expectRaw)
		got := normalizeOutput(gotRaw)
		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, err, input, expectRaw, gotRaw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

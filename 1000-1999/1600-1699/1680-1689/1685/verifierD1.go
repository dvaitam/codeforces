package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(4))
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(5) + 2 // 2..6
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		perm := r.Perm(n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(perm[i] + 1))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

// parsePermutation parses a space-separated line of integers into a slice.
func parsePermutation(s string) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	result := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		result[i] = v
	}
	return result, nil
}

// parseInputPermutation extracts n and p from the test input.
func parseInputPermutation(input string) (int, []int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 3 {
		return 0, nil, fmt.Errorf("expected at least 3 lines in input, got %d", len(lines))
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse n: %v", err)
	}
	p, err := parsePermutation(lines[2])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse p: %v", err)
	}
	if len(p) != n {
		return 0, nil, fmt.Errorf("expected %d elements in p, got %d", n, len(p))
	}
	return n, p, nil
}

// computeWeight computes the weight of permutation q given permutation p.
// weight = sum of |q[i] - p[q[i+1]]| for i=0..n-2, plus |q[n-1] - p[q[0]]|
// Both q and p are 1-indexed values.
func computeWeight(p, q []int) int {
	n := len(q)
	weight := 0
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		// q[next] is 1-indexed, p is 1-indexed
		weight += int(math.Abs(float64(q[i] - p[q[next]-1])))
	}
	return weight
}

// validatePermutation checks that q is a valid permutation of 1..n.
func validatePermutation(q []int, n int) error {
	if len(q) != n {
		return fmt.Errorf("expected %d elements, got %d", n, len(q))
	}
	seen := make([]bool, n+1)
	for i, v := range q {
		if v < 1 || v > n {
			return fmt.Errorf("element %d out of range [1,%d]: %d", i, n, v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate element: %d", v)
		}
		seen[v] = true
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD1 /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		fmt.Fprintln(os.Stderr, "REFERENCE_SOURCE_PATH not set")
		os.Exit(1)
	}
	ref := "refD1"
	cmd := exec.Command("go", "build", "-o", ref, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, input := range tests {
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary("./"+ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\n", i+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\n", i+1, rErr)
			os.Exit(1)
		}

		n, p, err := parseInputPermutation(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse input: %v\n", i+1, err)
			os.Exit(1)
		}

		// Parse and validate candidate output
		candQ, err := parsePermutation(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse candidate output: %v\noutput: %s", i+1, err, candOut)
			os.Exit(1)
		}
		if err := validatePermutation(candQ, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid candidate permutation: %v\noutput: %s", i+1, err, candOut)
			os.Exit(1)
		}

		// Parse reference output
		refQ, err := parsePermutation(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse reference output: %v\noutput: %s", i+1, err, refOut)
			os.Exit(1)
		}

		candWeight := computeWeight(p, candQ)
		refWeight := computeWeight(p, refQ)

		if candWeight != refWeight {
			fmt.Fprintf(os.Stderr, "case %d: candidate weight %d != reference weight %d\ninput:\n%scandidate: %s\nreference: %s",
				i+1, candWeight, refWeight, input, strings.TrimSpace(candOut), strings.TrimSpace(refOut))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

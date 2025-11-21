package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const maxN = 500000

func isPerfectSquare(x int64) bool {
	if x < 0 {
		return false
	}
	r := int64(math.Sqrt(float64(x)))
	for r*r < x {
		r++
	}
	for r*r > x {
		r--
	}
	return r*r == x
}

func sumIsSquare(n int) bool {
	sum := int64(n) * int64(n+1) / 2
	return isPerfectSquare(sum)
}

func buildTests() []int {
	tests := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 42, 100, 12345, maxN}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tests = append(tests, rng.Intn(maxN)+1)
	}
	return tests
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve binary path: %v\n", err)
		os.Exit(1)
	}

	tests := buildTests()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, n := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}

	output, err := runCandidate(bin, sb.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(strings.NewReader(output))
	for idx, n := range tests {
		impossible := sumIsSquare(n)
		var first int
		_, err := fmt.Fscan(reader, &first)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: expected value, got error %v\n", idx+1, n, err)
			os.Exit(1)
		}
		if impossible {
			if first != -1 {
				fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: expected -1 for impossible case, got %d\n", idx+1, n, first)
				os.Exit(1)
			}
			continue
		}
		perm := make([]int, n)
		perm[0] = first
		for i := 1; i < n; i++ {
			if _, err := fmt.Fscan(reader, &perm[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: reading permutation: %v\n", idx+1, n, err)
				os.Exit(1)
			}
		}
		if err := validatePermutation(perm); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (n=%d) failed: %v\n", idx+1, n, err)
			os.Exit(1)
		}
	}
	// ensure remaining tokens are just whitespace
	if _, err := fmt.Fscan(reader, new(int)); err == nil {
		fmt.Fprintln(os.Stderr, "extra output detected after processing all test cases")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func validatePermutation(perm []int) error {
	n := len(perm)
	seen := make([]bool, n+1)
	var sum int64
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d out of range [1,%d]", v, i+1, n)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
		sum += int64(v)
		if isPerfectSquare(sum) {
			return fmt.Errorf("prefix sum after %d elements (%d) is a perfect square", i+1, sum)
		}
	}
	for v := 1; v <= n; v++ {
		if !seen[v] {
			return fmt.Errorf("value %d missing in permutation", v)
		}
	}
	return nil
}

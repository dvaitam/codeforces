package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// generateCase creates a random test input and returns the input string and the array a.
// All a_i are even, as required by B1.
func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n/2+1) * 2 // even values in [0, n]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

// validateOutput semantically checks a candidate's output for problem 1970/B1.
// a is the input array (0-indexed).  output is the candidate's raw stdout.
// Returns nil if valid, or an error describing the violation.
func validateOutput(a []int, refVerdict string, output string) error {
	n := len(a)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	verdict := strings.TrimSpace(lines[0])
	if strings.EqualFold(verdict, "NO") {
		// Candidate says NO; the reference must also say NO.
		if strings.EqualFold(refVerdict, "YES") {
			return fmt.Errorf("candidate says NO but a valid placement exists")
		}
		return nil
	}

	if !strings.EqualFold(verdict, "YES") {
		return fmt.Errorf("first line must be YES or NO, got %q", verdict)
	}
	// Candidate says YES; reference must agree that YES is possible.
	if strings.EqualFold(refVerdict, "NO") {
		return fmt.Errorf("candidate says YES but reference says NO")
	}

	// We need n+1 more lines after the verdict line (n position lines + 1 direction line).
	if len(lines) < n+2 {
		return fmt.Errorf("expected %d lines after YES, got %d", n+1, len(lines)-1)
	}

	xs := make([]int, n)
	ys := make([]int, n)

	// Parse positions.
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[1+i])
		if len(parts) != 2 {
			return fmt.Errorf("position line %d: expected 2 integers, got %q", i+1, lines[1+i])
		}
		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("position line %d: parse error", i+1)
		}
		if x < 1 || x > n || y < 1 || y > n {
			return fmt.Errorf("wizard %d: position (%d,%d) out of range [1,%d]", i+1, x, y, n)
		}
		xs[i] = x
		ys[i] = y
	}

	// Check unique columns (all x_i distinct).
	colUsed := make(map[int]int) // column -> wizard (1-indexed)
	for i := 0; i < n; i++ {
		if prev, ok := colUsed[xs[i]]; ok {
			return fmt.Errorf("wizards %d and %d share column %d", prev, i+1, xs[i])
		}
		colUsed[xs[i]] = i + 1
	}

	// Parse direction line.
	dirParts := strings.Fields(lines[n+1])
	if len(dirParts) != n {
		return fmt.Errorf("direction line: expected %d integers, got %d", n, len(dirParts))
	}

	for i := 0; i < n; i++ {
		d, err := strconv.Atoi(dirParts[i])
		if err != nil {
			return fmt.Errorf("direction line: parse error at position %d", i+1)
		}
		if d < 1 || d > n {
			return fmt.Errorf("wizard %d: visit target %d out of range [1,%d]", i+1, d, n)
		}
		j := d - 1 // 0-indexed target
		dist := abs(xs[i]-xs[j]) + abs(ys[i]-ys[j])
		if dist != a[i] {
			return fmt.Errorf("wizard %d: distance to wizard %d is %d, expected %d", i+1, d, dist, a[i])
		}
	}

	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// solveB1 produces a reference verdict (YES or NO) for the given input.
func refVerdict(a []int) string {
	n := len(a)
	for i := 0; i < n; i++ {
		sel := a[i] / 2
		if i >= sel || i+sel < n {
			continue
		}
		return "NO"
	}
	return "YES"
}

func runCandidate(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB1 /path/to/candidate-binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 200; i++ {
		input, arr := generateCase(rng)
		rv := refVerdict(arr)
		output, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(arr, rv, output); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s", i+1, err, input, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

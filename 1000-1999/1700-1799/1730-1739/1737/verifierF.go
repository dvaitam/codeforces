package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return out.String() + errb.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int, []int) {
	m := rng.Intn(3) + 1
	exps := make([]int, m)
	for i := 0; i < m; i++ {
		exps[i] = rng.Intn(3) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", m)
	for i, v := range exps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), m, exps
}

// Check whether the output is a valid answer for the given exponents.
// Returns nil if valid, error otherwise.
func checkAnswer(output string, m int, b []int) error {
	N := 1
	for _, v := range b {
		N *= (v + 1)
	}
	// expected number of output rows = N - 1 (all divisors > 1)
	expectedRows := N - 1

	output = strings.TrimSpace(output)
	if output == "-1" {
		// Candidate says impossible; verify that it truly is impossible.
		// Conditions for impossibility (from the problem/reference):
		cnt := [5]int{}
		for _, v := range b {
			c := v
			if c > 4 {
				c = 4
			}
			cnt[c]++
		}
		impossible := cnt[4] > 0 || cnt[3] >= 2 || (cnt[3] >= 1 && cnt[2] >= 1) || cnt[2] >= 3
		if !impossible {
			return fmt.Errorf("candidate output -1 but a solution exists")
		}
		return nil
	}

	lines := strings.Split(output, "\n")
	if len(lines) != expectedRows {
		return fmt.Errorf("expected %d rows, got %d", expectedRows, len(lines))
	}

	type tuple struct {
		vals string
	}
	seen := make(map[string]bool)

	prev := make([]int, m)
	for rowIdx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != m {
			return fmt.Errorf("row %d: expected %d values, got %d", rowIdx+1, m, len(fields))
		}
		cur := make([]int, m)
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("row %d col %d: invalid integer %q", rowIdx+1, j+1, f)
			}
			if v < 0 || v > b[j] {
				return fmt.Errorf("row %d col %d: value %d out of range [0, %d]", rowIdx+1, j+1, v, b[j])
			}
			cur[j] = v
		}

		// Check > 1 (at least one exponent must be positive)
		allZero := true
		for _, v := range cur {
			if v > 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return fmt.Errorf("row %d: all exponents are zero (represents 1)", rowIdx+1)
		}

		// Check distinct
		key := fmt.Sprint(cur)
		if seen[key] {
			return fmt.Errorf("row %d: duplicate entry %v", rowIdx+1, cur)
		}
		seen[key] = true

		// Check consecutive GCD is prime:
		// Two divisors' GCD is prime iff they differ in exactly one coordinate by exactly 1.
		if rowIdx > 0 {
			diffCount := 0
			for j := 0; j < m; j++ {
				d := cur[j] - prev[j]
				if d < 0 {
					d = -d
				}
				if d == 1 {
					diffCount++
				} else if d > 1 {
					return fmt.Errorf("row %d: coordinate %d differs by %d from previous (need exactly 1 diff of 1)", rowIdx+1, j+1, d)
				}
			}
			if diffCount != 1 {
				return fmt.Errorf("row %d: %d coordinates differ by 1 from previous (need exactly 1)", rowIdx+1, diffCount)
			}
		}
		copy(prev, cur)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Run tests with various random cases
	for i := 0; i < 100; i++ {
		input, m, exps := genCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkAnswer(got, m, exps); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}

	// Also test known edge cases
	edgeCases := []struct {
		input string
		m     int
		b     []int
	}{
		{"1\n1\n1\n", 1, []int{1}},
		{"1\n1\n3\n", 1, []int{3}},
		{"1\n1\n4\n", 1, []int{4}},
		{"1\n2\n1 1\n", 2, []int{1, 1}},
		{"1\n2\n2 2\n", 2, []int{2, 2}},
		{"1\n3\n1 1 1\n", 3, []int{1, 1, 1}},
	}
	for i, ec := range edgeCases {
		got, err := runProg(bin, ec.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on edge case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkAnswer(got, ec.m, ec.b); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on edge case %d: %v\ninput:\n%sgot:\n%s\n", i+1, err, ec.input, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}

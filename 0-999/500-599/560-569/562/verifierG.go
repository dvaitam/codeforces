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

type Rule struct{ a, b, c int }

type TestCase struct {
	n     int
	rules []Rule
}

func (tc TestCase) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for _, r := range tc.rules {
		fmt.Fprintf(&sb, "%d %d %d\n", r.a, r.b, r.c)
	}
	return sb.String()
}

// genCase generates a valid input: each server 1..n appears exactly 4 times
// in {a_i} and exactly 8 times in {b_i, c_i}.
func genCase(rng *rand.Rand) TestCase {
	n := rng.Intn(3) + 1
	total := 4 * n

	aVals := make([]int, total)
	for i := range aVals {
		aVals[i] = i/4 + 1
	}
	rng.Shuffle(len(aVals), func(i, j int) { aVals[i], aVals[j] = aVals[j], aVals[i] })

	bcVals := make([]int, 8*n)
	for i := range bcVals {
		bcVals[i] = i/8 + 1
	}
	rng.Shuffle(len(bcVals), func(i, j int) { bcVals[i], bcVals[j] = bcVals[j], bcVals[i] })

	rules := make([]Rule, total)
	for i := range rules {
		rules[i] = Rule{aVals[i], bcVals[2*i], bcVals[2*i+1]}
	}
	return TestCase{n, rules}
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// validate checks that the YES answer is a valid ordering (permutation of 1..4n
// that never exceeds 9 processes on any server at any step).
// Initially each server has 4 processes. Rule i: destroy one on a[i], create one
// on b[i], create one on c[i]. The net delta is tracked in s[]; total = 4 + s[k].
// Constraint: 4 + s[k] <= 9 → s[k] <= 5, checked only when s increases (b/c).
func validate(tc TestCase, output string) error {
	total := 4 * tc.n
	lines := strings.SplitN(strings.TrimSpace(output), "\n", 2)
	if strings.TrimSpace(lines[0]) == "NO" {
		return fmt.Errorf("got NO (for valid inputs the answer is always YES)")
	}
	if strings.TrimSpace(lines[0]) != "YES" {
		return fmt.Errorf("first line must be YES or NO, got %q", lines[0])
	}
	if len(lines) < 2 {
		return fmt.Errorf("missing permutation line")
	}
	fields := strings.Fields(lines[1])
	if len(fields) != total {
		return fmt.Errorf("expected %d numbers, got %d", total, len(fields))
	}

	s := make([]int, tc.n+1)
	seen := make([]bool, total+1)
	for step, f := range fields {
		rule, err := strconv.Atoi(f)
		if err != nil || rule < 1 || rule > total {
			return fmt.Errorf("invalid rule number %q at position %d", f, step+1)
		}
		if seen[rule] {
			return fmt.Errorf("duplicate rule %d at position %d", rule, step+1)
		}
		seen[rule] = true
		r := tc.rules[rule-1]
		s[r.a]--
		s[r.b]++
		s[r.c]++
		if s[r.b] > 5 {
			return fmt.Errorf("server %d has %d processes at step %d (exceeds 9)", r.b, 4+s[r.b], step+1)
		}
		if s[r.c] > 5 {
			return fmt.Errorf("server %d has %d processes at step %d (exceeds 9)", r.c, 4+s[r.c], step+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Sample from problem statement (n=1: all rules are 1->1,1)
	samples := []TestCase{
		{1, []Rule{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}, {1, 1, 1}}},
	}

	allCases := append(samples, func() []TestCase {
		tc := make([]TestCase, 100)
		for i := range tc {
			tc[i] = genCase(rng)
		}
		return tc
	}()...)

	for i, tc := range allCases {
		input := tc.Input()
		got, err := runBinary(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validate(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(allCases))
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Case struct {
	n, m int
	a    []int
	raw  string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []Case

	addRaw := func(n, m int, a []int) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[j])
		}
		sb.WriteByte('\n')
		cases = append(cases, Case{n, m, a, sb.String()})
	}

	// Edge cases.
	addRaw(2, 2, []int{1, 2})
	addRaw(2, 2, []int{1, 1})
	addRaw(2, 1, []int{3, 5})     // m=1, only odd available -> -1
	addRaw(4, 10, []int{2, 4, 6, 8}) // all even
	addRaw(4, 10, []int{1, 3, 5, 7}) // all odd
	addRaw(4, 4, []int{1, 1, 1, 1})
	addRaw(2, 1000000000, []int{999999999, 1000000000})

	// Random cases.
	for len(cases) < 100 {
		n := (rng.Intn(10) + 1) * 2
		m := rng.Intn(200) + n
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(m) + 1
		}
		addRaw(n, m, a)
	}
	return cases
}

// computeMinExchanges computes the minimum number of exchanges needed.
// Returns -1 if impossible.
func computeMinExchanges(n, m int, a []int) int {
	target := n / 2 // need this many even AND this many odd

	// Count available even/odd values in [1..m].
	totalEven := m / 2
	totalOdd := (m + 1) / 2

	// First pass: keep distinct values, count even/odd among kept.
	used := make(map[int]bool)
	keepEven, keepOdd := 0, 0
	dupCount := 0
	isDup := make([]bool, n)
	for i, v := range a {
		if used[v] {
			isDup[i] = true
			dupCount++
		} else {
			used[v] = true
			if v%2 == 0 {
				keepEven++
			} else {
				keepOdd++
			}
		}
	}

	// Must replace all duplicates. Additionally may need to convert
	// some kept values if we have too many even or too many odd.
	changes := dupCount
	extraEven := keepEven - target // positive = too many even kept
	extraOdd := keepOdd - target   // positive = too many odd kept

	if extraEven > 0 {
		changes += extraEven
	}
	if extraOdd > 0 {
		changes += extraOdd
	}

	// After changes, count how many even/odd we need from [1..m].
	needEven := target - keepEven
	if needEven < 0 {
		needEven = 0
	}
	needOdd := target - keepOdd
	if needOdd < 0 {
		needOdd = 0
	}

	// Available = total in [1..m] minus those already kept.
	keptInRange := 0
	for v := range used {
		if v >= 1 && v <= m {
			keptInRange++ // these are "used" from Nikolay's perspective
		}
	}
	// Actually we need to count kept even in [1..m] and kept odd in [1..m].
	keptEvenInRange, keptOddInRange := 0, 0
	for v := range used {
		if v >= 1 && v <= m {
			if v%2 == 0 {
				keptEvenInRange++
			} else {
				keptOddInRange++
			}
		}
	}
	// But after removing extra kept values, those become available again.
	availEven := totalEven - keptEvenInRange
	availOdd := totalOdd - keptOddInRange
	if extraEven > 0 {
		// We'll un-keep some even values; if they're in [1..m] they become available.
		// Conservative: just add extraEven back (they might be > m, but that's ok, we check later).
		availEven += extraEven
	}
	if extraOdd > 0 {
		availOdd += extraOdd
	}

	if availEven < needEven || availOdd < needOdd {
		return -1
	}
	return changes
}

// validate checks candidate output for a given test case.
func validate(c Case, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "-1" {
		// Candidate says impossible. Check if that's correct.
		minEx := computeMinExchanges(c.n, c.m, c.a)
		if minEx == -1 {
			return nil
		}
		return fmt.Errorf("candidate says -1 but solution exists with %d exchanges", minEx)
	}
	if len(lines) < 2 {
		return fmt.Errorf("expected 2 lines of output, got %d", len(lines))
	}

	kStr := strings.TrimSpace(lines[0])
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid exchange count: %q", kStr)
	}

	tokens := strings.Fields(lines[1])
	if len(tokens) != c.n {
		return fmt.Errorf("expected %d values in result, got %d", c.n, len(tokens))
	}
	res := make([]int, c.n)
	for i, tok := range tokens {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("invalid integer at position %d: %q", i+1, tok)
		}
		res[i] = v
	}

	// Count actual exchanges.
	actualChanges := 0
	for i := 0; i < c.n; i++ {
		if res[i] != c.a[i] {
			actualChanges++
		}
	}
	if actualChanges != k {
		return fmt.Errorf("claimed %d exchanges but %d positions differ", k, actualChanges)
	}

	// Check all values distinct.
	seen := make(map[int]bool, c.n)
	for i, v := range res {
		if seen[v] {
			return fmt.Errorf("duplicate value %d at position %d", v, i+1)
		}
		seen[v] = true
	}

	// Check even/odd balance.
	evenCount := 0
	for _, v := range res {
		if v%2 == 0 {
			evenCount++
		}
	}
	target := c.n / 2
	if evenCount != target {
		return fmt.Errorf("expected %d even values, got %d", target, evenCount)
	}

	// Check exchanged cards are valid (value in [1..m] and positive).
	for i := 0; i < c.n; i++ {
		if res[i] != c.a[i] {
			if res[i] < 1 || res[i] > c.m {
				return fmt.Errorf("exchanged card at position %d has value %d outside [1, %d]", i+1, res[i], c.m)
			}
		}
	}

	// Check minimality.
	minEx := computeMinExchanges(c.n, c.m, c.a)
	if minEx == -1 {
		return fmt.Errorf("no solution exists but candidate produced one")
	}
	if k != minEx {
		return fmt.Errorf("not minimal: candidate used %d exchanges, minimum is %d", k, minEx)
	}

	return nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		got, err := runBinary(bin, c.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.raw)
			os.Exit(1)
		}
		if err := validate(c, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ncandidate output:\n%s\ninput:\n%s", i+1, err, got, c.raw)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

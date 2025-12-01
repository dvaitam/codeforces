package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// refSourceA3 points to the local reference so Go does not try to resolve a
	// module-like path under GOPATH.
	refSourceA3 = "207A3.go"
	orderLimit  = 200000
)

type scientist struct {
	k        int
	a1, x, y int64
	m        int64
	seq      []int64
}

type testCase struct {
	scis []scientist
}

type pair struct {
	val int64
	id  int
}

type parsedOutput struct {
	bad   int64
	order []pair
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA3.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := buildInput(tc)
		total := tc.totalProblems()
		needOrder := total <= orderLimit

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}
		refRes, err := parseOutput(refOut, total, needOrder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candRes, err := parseOutput(candOut, total, needOrder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse candidate output on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}

		if candRes.bad != refRes.bad {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d bad pairs got %d\n", idx+1, refRes.bad, candRes.bad)
			os.Exit(1)
		}

		if len(candRes.order) > 0 {
			if err := validateOrder(tc, candRes.order); err != nil {
				fmt.Fprintf(os.Stderr, "invalid order on test %d: %v\n", idx+1, err)
				os.Exit(1)
			}
			if countBadPairs(candRes.order) != candRes.bad {
				fmt.Fprintf(os.Stderr, "bad pair count mismatch on test %d\n", idx+1)
				os.Exit(1)
			}
		} else if needOrder {
			fmt.Fprintf(os.Stderr, "missing order for test %d with %d problems\n", idx+1, total)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "207A3-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceA3))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, total int, requireOrder bool) (parsedOutput, error) {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return parsedOutput{}, fmt.Errorf("empty output")
	}
	bad, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return parsedOutput{}, fmt.Errorf("invalid bad pair count")
	}
	tokens = tokens[1:]

	expectOrder := requireOrder
	if !expectOrder && len(tokens) > 0 {
		if len(tokens) != 2*total {
			return parsedOutput{}, fmt.Errorf("unexpected extra tokens in output")
		}
		expectOrder = true
	}

	res := parsedOutput{bad: bad}
	if expectOrder {
		if len(tokens) != 2*total {
			return parsedOutput{}, fmt.Errorf("expected %d order entries, got %d tokens", total, len(tokens)/2)
		}
		order := make([]pair, 0, total)
		for i := 0; i < len(tokens); i += 2 {
			val, err1 := strconv.ParseInt(tokens[i], 10, 64)
			id, err2 := strconv.Atoi(tokens[i+1])
			if err1 != nil || err2 != nil {
				return parsedOutput{}, fmt.Errorf("invalid order entry at position %d", i/2+1)
			}
			order = append(order, pair{val: val, id: id})
		}
		res.order = order
	}
	return res, nil
}

func validateOrder(tc testCase, order []pair) error {
	total := tc.totalProblems()
	if len(order) != total {
		return fmt.Errorf("order length mismatch: got %d expected %d", len(order), total)
	}
	positions := make([]int, len(tc.scis))
	for idx, entry := range order {
		if entry.id < 1 || entry.id > len(tc.scis) {
			return fmt.Errorf("invalid scientist id %d at position %d", entry.id, idx+1)
		}
		s := tc.scis[entry.id-1]
		pos := positions[entry.id-1]
		if pos >= s.k {
			return fmt.Errorf("too many tasks taken from scientist %d", entry.id)
		}
		if s.seq[pos] != entry.val {
			return fmt.Errorf("value mismatch for scientist %d at task %d", entry.id, pos+1)
		}
		positions[entry.id-1]++
	}
	for i, pos := range positions {
		if pos != tc.scis[i].k {
			return fmt.Errorf("not all tasks for scientist %d were used", i+1)
		}
	}
	return nil
}

func countBadPairs(order []pair) int64 {
	if len(order) == 0 {
		return 0
	}
	var bad int64
	prev := order[0].val
	for i := 1; i < len(order); i++ {
		if order[i].val < prev {
			bad++
		}
		prev = order[i].val
	}
	return bad
}

func buildInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tc.scis))
	for _, s := range tc.scis {
		fmt.Fprintf(&b, "%d %d %d %d %d\n", s.k, s.a1, s.x, s.y, s.m)
	}
	return b.String()
}

func (tc testCase) totalProblems() int {
	total := 0
	for _, s := range tc.scis {
		total += s.k
	}
	return total
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, sampleCase1())
	tests = append(tests, sampleCase2())
	tests = append(tests, equalValuesCase())
	tests = append(tests, singleScientistCase())
	tests = append(tests, alternatingPressureCase())
	tests = append(tests, denseLimitCase())
	tests = append(tests, bigNoOrderCase())

	rng := rand.New(rand.NewSource(20720207))
	for i := 0; i < 18; i++ {
		target := rng.Intn(orderLimit-1) + 1
		tests = append(tests, randomCase(rng, target, 60))
	}
	for i := 0; i < 8; i++ {
		target := orderLimit + rng.Intn(60000) + 10
		tests = append(tests, randomCase(rng, target, 200))
	}
	return tests
}

func sampleCase1() testCase {
	scis := []scientist{
		buildScientist(2, 1, 1, 1, 10),
		buildScientist(2, 3, 1, 1, 10),
	}
	return testCase{scis: scis}
}

func sampleCase2() testCase {
	scis := []scientist{
		buildScientist(3, 10, 2, 3, 1000),
		buildScientist(3, 100, 1, 999, 1000),
	}
	return testCase{scis: scis}
}

func equalValuesCase() testCase {
	scis := []scientist{
		buildScientist(4, 5, 1, 6, 6),
		buildScientist(4, 5, 1, 6, 6),
		buildScientist(4, 5, 1, 6, 6),
	}
	return testCase{scis: scis}
}

func singleScientistCase() testCase {
	return testCase{scis: []scientist{buildScientist(7, 2, 3, 5, 1000)}}
}

func alternatingPressureCase() testCase {
	scis := []scientist{
		buildScientist(6, 1000, 999983, 7, 1_000_000_000),
		buildScientist(6, 1, 1, 1, 2),
	}
	return testCase{scis: scis}
}

func denseLimitCase() testCase {
	const perScientist = 40
	need := orderLimit
	n := need / perScientist
	scis := make([]scientist, n)
	current := int64(0)
	for i := 0; i < n; i++ {
		scis[i] = buildScientist(perScientist, current%1_000_000_000, 1, 3, 1_000_000_000)
		current += 7
	}
	return testCase{scis: scis}
}

func bigNoOrderCase() testCase {
	target := orderLimit + 50000
	rng := rand.New(rand.NewSource(1337))
	return randomCase(rng, target, 500)
}

func randomCase(rng *rand.Rand, targetTotal, maxN int) testCase {
	if targetTotal <= 0 {
		targetTotal = 1
	}
	maxN = minInt(maxN, targetTotal)
	maxN = minInt(maxN, 5000)
	minN := (targetTotal + 5000 - 1) / 5000
	if minN < 1 {
		minN = 1
	}
	if maxN < minN {
		maxN = minN
	}
	n := minN
	if maxN > minN {
		n = minN + rng.Intn(maxN-minN+1)
	}
	ks := distributeCounts(rng, targetTotal, n)
	scis := make([]scientist, n)
	for i, k := range ks {
		scis[i] = randomScientist(rng, k)
	}
	return testCase{scis: scis}
}

func distributeCounts(rng *rand.Rand, total, n int) []int {
	ks := make([]int, n)
	remaining := total
	for i := 0; i < n; i++ {
		left := n - i
		minVal := maxInt(1, remaining-(left-1)*5000)
		maxVal := minInt(5000, remaining-(left-1))
		if minVal > maxVal {
			minVal = maxVal
		}
		k := minVal
		if maxVal > minVal {
			k = minVal + rng.Intn(maxVal-minVal+1)
		}
		ks[i] = k
		remaining -= k
	}
	return ks
}

func randomScientist(rng *rand.Rand, k int) scientist {
	m := rng.Int63n(1_000_000_000) + 1
	if m == 0 {
		m = 1
	}
	a1 := rng.Int63n(m)
	x := rng.Int63n(1_000_000_000) + 1
	y := rng.Int63n(1_000_000_000) + 1
	return buildScientist(k, a1, x, y, m)
}

func buildScientist(k int, a1, x, y, m int64) scientist {
	return scientist{
		k:   k,
		a1:  a1,
		x:   x,
		y:   y,
		m:   m,
		seq: generateSequence(k, a1, x, y, m),
	}
}

func generateSequence(k int, a1, x, y, m int64) []int64 {
	seq := make([]int64, k)
	if k == 0 {
		return seq
	}
	seq[0] = a1
	prev := a1
	for i := 1; i < k; i++ {
		prev = (prev*x + y) % m
		seq[i] = prev
	}
	return seq
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

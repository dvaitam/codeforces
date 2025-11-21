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
	refSourceA1 = "0-999/200-299/200-209/207/207A1.go"
	orderLimit  = 200000
)

type pair struct {
	val int64
	id  int
}

type scientist struct {
	k       int
	a1, x   int64
	y, m    int64
	seq     []int64
	hasSeqs bool
}

type testCase struct {
	scientists []scientist
	total      int
}

func (tc testCase) needsOrder() bool {
	return tc.total <= orderLimit
}

func (tc testCase) buildInput() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tc.scientists))
	for _, s := range tc.scientists {
		fmt.Fprintf(&b, "%d %d %d %d %d\n", s.k, s.a1, s.x, s.y, s.m)
	}
	return b.String()
}

func (tc testCase) validateOrder(order []pair) error {
	if !tc.needsOrder() {
		return nil
	}
	if len(order) != tc.total {
		return fmt.Errorf("expected %d lines in order, got %d", tc.total, len(order))
	}
	pos := make([]int, len(tc.scientists))
	for _, p := range order {
		if p.id < 1 || p.id > len(tc.scientists) {
			return fmt.Errorf("invalid scientist id %d", p.id)
		}
		sci := tc.scientists[p.id-1]
		if !sci.hasSeqs {
			return fmt.Errorf("missing precomputed sequence for scientist %d", p.id)
		}
		idx := pos[p.id-1]
		if idx >= sci.k {
			return fmt.Errorf("scientist %d tasks exhausted", p.id)
		}
		if sci.seq[idx] != p.val {
			return fmt.Errorf("value mismatch for scientist %d at task %d: expected %d got %d", p.id, idx+1, sci.seq[idx], p.val)
		}
		pos[p.id-1]++
	}
	for i, used := range pos {
		if used != tc.scientists[i].k {
			return fmt.Errorf("scientist %d tasks not fully scheduled", i+1)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		input := tc.buildInput()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, refOut)
			os.Exit(1)
		}
		refBad, err := parseReferenceBadPairs(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, candOut)
			os.Exit(1)
		}
		candBad, order, err := parseCandidateOutput(candOut, tc.needsOrder(), tc.total)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, candOut)
			os.Exit(1)
		}
		if candBad != refBad {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d bad pairs got %d\n", i+1, refBad, candBad)
			os.Exit(1)
		}
		if tc.needsOrder() {
			if err := tc.validateOrder(order); err != nil {
				fmt.Fprintf(os.Stderr, "order check failed on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			if countBadPairs(order) != candBad {
				fmt.Fprintf(os.Stderr, "bad pair recount mismatch on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "207A1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceA1))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseReferenceBadPairs(out string) (int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	return strconv.ParseInt(tokens[0], 10, 64)
}

func parseCandidateOutput(out string, needOrder bool, total int) (int64, []pair, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	bad, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid bad pair count")
	}
	tokens = tokens[1:]
	if needOrder {
		expected := total * 2
		if len(tokens) != expected {
			return 0, nil, fmt.Errorf("expected %d integers for the order, got %d", expected, len(tokens))
		}
		order := make([]pair, total)
		for i := 0; i < total; i++ {
			val, err1 := strconv.ParseInt(tokens[2*i], 10, 64)
			id, err2 := strconv.Atoi(tokens[2*i+1])
			if err1 != nil || err2 != nil {
				return 0, nil, fmt.Errorf("invalid order entry at position %d", i+1)
			}
			order[i] = pair{val: val, id: id}
		}
		return bad, order, nil
	}
	if len(tokens) != 0 {
		return 0, nil, fmt.Errorf("unexpected extra output when order is not required")
	}
	return bad, nil, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, sampleTest1())
	tests = append(tests, sampleTest2())
	tests = append(tests, zeroScientistTasks())
	tests = append(tests, manyScientistsSingleTask())
	tests = append(tests, boundaryOrderCase())
	tests = append(tests, largeNoOrderCase())

	rng := rand.New(rand.NewSource(207001))
	for len(tests) < 25 {
		tests = append(tests, randomCase(rng, rng.Intn(6)+1, rng.Intn(50)+1))
	}
	for len(tests) < 40 {
		tests = append(tests, randomCase(rng, rng.Intn(30)+5, rng.Intn(200)+1))
	}

	return tests
}

func sampleTest1() testCase {
	// Matches statement sample 1.
	scis := []scientist{
		{k: 2, a1: 1, x: 1, y: 1, m: 10},
		{k: 2, a1: 3, x: 1, y: 1, m: 10},
	}
	return buildTestCase(scis)
}

func sampleTest2() testCase {
	scis := []scientist{
		{k: 3, a1: 10, x: 2, y: 3, m: 1000},
		{k: 3, a1: 100, x: 1, y: 999, m: 1000},
	}
	return buildTestCase(scis)
}

func zeroScientistTasks() testCase {
	scis := []scientist{
		{k: 1, a1: 0, x: 1, y: 1, m: 2},
		{k: 5, a1: 7, x: 3, y: 2, m: 50},
	}
	return buildTestCase(scis)
}

func manyScientistsSingleTask() testCase {
	n := 5000
	scis := make([]scientist, n)
	for i := 0; i < n; i++ {
		val := int64(i % 7)
		scis[i] = scientist{k: 1, a1: val, x: 1, y: 1, m: 999_999_937}
	}
	return buildTestCase(scis)
}

func largeNoOrderCase() testCase {
	n := 60
	scis := make([]scientist, n)
	for i := 0; i < n; i++ {
		k := 3500
		if i%3 == 1 {
			k = 2800
		} else if i%3 == 2 {
			k = 4200
		}
		m := int64(900_000_000 - i*1000)
		scis[i] = scientist{
			k:  k,
			a1: int64((i * 37) % int(m)),
			x:  999_999_937,
			y:  34567,
			m:  m,
		}
	}
	return buildTestCase(scis)
}

func boundaryOrderCase() testCase {
	n := 100
	scis := make([]scientist, n)
	for i := 0; i < n; i++ {
		k := 2000
		m := int64(800_000_000 + i*1000)
		scis[i] = scientist{
			k:  k,
			a1: int64((i * 73) % int(m)),
			x:  int64(1000000000 - i*1234),
			y:  int64(500000000 - i*4321),
			m:  m,
		}
		if scis[i].x <= 0 {
			scis[i].x = 1
		}
		if scis[i].y <= 0 {
			scis[i].y = 1
		}
	}
	return buildTestCase(scis)
}

func randomCase(rng *rand.Rand, n, maxK int) testCase {
	scis := make([]scientist, n)
	for i := 0; i < n; i++ {
		k := rng.Intn(maxK) + 1
		m := int64(rng.Intn(1_000_000_000) + 1)
		a1 := int64(rng.Int63n(m))
		x := int64(rng.Intn(1_000_000_000) + 1)
		y := int64(rng.Intn(1_000_000_000) + 1)
		scis[i] = scientist{
			k: k, a1: a1, x: x, y: y, m: m,
		}
	}
	return buildTestCase(scis)
}

func buildTestCase(scis []scientist) testCase {
	tc := testCase{
		scientists: make([]scientist, len(scis)),
	}
	total := 0
	for i, s := range scis {
		tc.scientists[i] = s
		total += s.k
	}
	tc.total = total
	if tc.needsOrder() {
		for i := range tc.scientists {
			tc.scientists[i].seq = generateSequence(tc.scientists[i])
			tc.scientists[i].hasSeqs = true
		}
	}
	return tc
}

func generateSequence(s scientist) []int64 {
	if s.k == 0 {
		return nil
	}
	seq := make([]int64, s.k)
	seq[0] = s.a1
	for i := 1; i < s.k; i++ {
		seq[i] = (seq[i-1]*s.x + s.y) % s.m
	}
	return seq
}

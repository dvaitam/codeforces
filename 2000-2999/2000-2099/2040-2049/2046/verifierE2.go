package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type participant struct {
	a int64
	b int64
	s int64
}

type cityTest struct {
	name   string
	parts  []participant
	cities [][]int
}

type outputCase struct {
	possible bool
	probs    []problem
}

type problem struct {
	d int64
	t int64
}

const (
	refSource  = "2046E2.go"
	limitVal   = int64(1_000_000_000)
	maxPerPart = 5
)

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidatePath := os.Args[len(os.Args)-1]

	refBin, refCleanup, err := buildBinary(referencePath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := prepareCandidate(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}
	refParsed, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}
	for idx, sol := range refParsed {
		if err := validateSolution(sol, tests[idx], sol.possible); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid solution on test %d (%s): %v\n", idx+1, tests[idx].name, err)
			os.Exit(1)
		}
	}

	candOut, err := runProgram(candBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}
	candParsed, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\n%s", err, candOut)
		fmt.Fprintln(os.Stderr, previewInput(input))
		os.Exit(1)
	}

	if len(candParsed) != len(refParsed) {
		fmt.Fprintf(os.Stderr, "mismatched test case count: reference %d candidate %d\n", len(refParsed), len(candParsed))
		os.Exit(1)
	}

	for i := range tests {
		expectPossible := refParsed[i].possible
		if err := validateSolution(candParsed[i], tests[i], expectPossible); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\n", i+1, tests[i].name, err)
			fmt.Fprintln(os.Stderr, previewInput(input))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), refSource)
	}
	return refSource
}

func buildBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "ref-2046E2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "refbin")
	cmd := exec.Command("go", "build", "-o", bin, filepath.Base(src))
	cmd.Dir = filepath.Dir(src)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "cand-2046E2-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", filepath.Base(path))
		cmd.Dir = filepath.Dir(path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(out string, tests []cityTest) ([]outputCase, error) {
	tokens := strings.Fields(out)
	res := make([]outputCase, 0, len(tests))
	idx := 0
	for ti, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("not enough tokens for test %d", ti+1)
		}
		if tokens[idx] == "-1" {
			res = append(res, outputCase{possible: false})
			idx++
			continue
		}
		p, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid integer p on test %d: %v", ti+1, err)
		}
		if p < 1 || p > maxPerPart*tcSize(tc) {
			return nil, fmt.Errorf("invalid p=%d on test %d", p, ti+1)
		}
		idx++
		if idx+2*p-1 >= len(tokens) {
			return nil, fmt.Errorf("not enough tokens for problems on test %d", ti+1)
		}
		probs := make([]problem, p)
		for i := 0; i < p; i++ {
			dVal, err := strconv.ParseInt(tokens[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid difficulty on test %d problem %d: %v", ti+1, i+1, err)
			}
			tVal, err := strconv.ParseInt(tokens[idx+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid topic on test %d problem %d: %v", ti+1, i+1, err)
			}
			probs[i] = problem{d: dVal, t: tVal}
			idx += 2
		}
		res = append(res, outputCase{possible: true, probs: probs})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens at end of output")
	}
	return res, nil
}

func validateSolution(sol outputCase, tc cityTest, expectPossible bool) error {
	n := tcSize(tc)
	if !expectPossible {
		if sol.possible {
			return fmt.Errorf("expected -1 but got a solution")
		}
		return nil
	}
	if !sol.possible {
		return fmt.Errorf("expected a solution but got -1")
	}
	if len(sol.probs) == 0 || len(sol.probs) > maxPerPart*n {
		return fmt.Errorf("problem count %d out of bounds", len(sol.probs))
	}
	topicSeen := make(map[int64]bool)
	for i, pr := range sol.probs {
		if pr.d < 0 || pr.d > limitVal || pr.t < 0 || pr.t > limitVal {
			return fmt.Errorf("problem %d has value out of range", i+1)
		}
		if topicSeen[pr.t] {
			return fmt.Errorf("duplicate topic %d", pr.t)
		}
		topicSeen[pr.t] = true
	}

	// Count solved problems per participant.
	solved := make([]int, n)
	for _, pr := range sol.probs {
		for idx, p := range tc.parts {
			if p.a >= pr.d || (p.s == pr.t && p.b >= pr.d) {
				solved[idx]++
			}
		}
	}

	minCity := make([]int, len(tc.cities))
	maxCity := make([]int, len(tc.cities))
	for i := range tc.cities {
		minCity[i] = 1 << 60
		maxCity[i] = -1
		for _, id := range tc.cities[i] {
			val := solved[id]
			if val < minCity[i] {
				minCity[i] = val
			}
			if val > maxCity[i] {
				maxCity[i] = val
			}
		}
	}

	for i := 0; i < len(tc.cities); i++ {
		for j := i + 1; j < len(tc.cities); j++ {
			if minCity[i] <= maxCity[j] {
				return fmt.Errorf("city %d does not dominate city %d", i+1, j+1)
			}
		}
	}

	return nil
}

func tcSize(tc cityTest) int {
	return len(tc.parts)
}

func buildTests() []cityTest {
	tests := make([]cityTest, 0, 200)
	tests = append(tests, deterministicTests()...)
	tests = append(tests, randomTests(150)...)
	return tests
}

func deterministicTests() []cityTest {
	return []cityTest{
		{
			name: "simple-possible",
			parts: []participant{
				{a: 5, b: 5, s: 1},
				{a: 1, b: 1, s: 2},
			},
			cities: [][]int{{0}, {1}},
		},
		{
			name: "identical-impossible",
			parts: []participant{
				{a: 1, b: 1, s: 1},
				{a: 1, b: 1, s: 1},
			},
			cities: [][]int{{0}, {1}},
		},
		{
			name: "three-cities",
			parts: []participant{
				{a: 3, b: 5, s: 1},
				{a: 2, b: 2, s: 2},
				{a: 1, b: 4, s: 1},
				{a: 0, b: 3, s: 3},
			},
			cities: [][]int{{0, 1}, {2}, {3}},
		},
	}
}

func randomTests(count int) []cityTest {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]cityTest, 0, count)
	for len(tests) < count {
		n := rng.Intn(60) + 2
		m := rng.Intn(n-1) + 2
		parts := make([]participant, n)
		for i := 0; i < n; i++ {
			a := rng.Int63n(30)
			b := a + rng.Int63n(30)
			s := rng.Int63n(50)
			parts[i] = participant{a: a, b: b, s: s}
		}

		cities := make([][]int, m)
		order := rng.Perm(n)
		for i := 0; i < m; i++ {
			cities[i] = []int{order[i]}
		}
		for i := m; i < n; i++ {
			c := rng.Intn(m)
			cities[c] = append(cities[c], order[i])
		}

		tests = append(tests, cityTest{
			name:   fmt.Sprintf("rand-%d", len(tests)+1),
			parts:  parts,
			cities: cities,
		})
	}
	return tests
}

func buildInput(tests []cityTest) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 128)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		n := len(tc.parts)
		m := len(tc.cities)
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		for _, p := range tc.parts {
			sb.WriteString(strconv.FormatInt(p.a, 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(p.b, 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(p.s, 10))
			sb.WriteByte('\n')
		}
		for _, city := range tc.cities {
			sb.WriteString(strconv.Itoa(len(city)))
			for _, idx := range city {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(idx + 1))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func previewInput(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 25 {
		lines = lines[:25]
	}
	return strings.Join(lines, "\n")
}

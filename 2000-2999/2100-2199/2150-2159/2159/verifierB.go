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

type testCase struct {
	name  string
	input string
	count int
}

type gridCase struct {
	rows []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2159B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2159B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(cases []gridCase) (string, int) {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	total := 0
	for _, cs := range cases {
		n := len(cs.rows)
		m := len(cs.rows[0])
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		for _, row := range cs.rows {
			sb.WriteString(row)
			sb.WriteByte('\n')
			total += len(row)
		}
	}
	return sb.String(), total
}

func deterministicTests() []testCase {
	tests := []testCase{}

	// All ones 2x2 grid -> every cell has rectangle area 4.
	caseAllOnes := []gridCase{
		{rows: []string{
			"11",
			"11",
		}},
	}
	if input, count := formatInput(caseAllOnes); count > 0 {
		tests = append(tests, testCase{name: "all_ones_small", input: input, count: count})
	}

	// Grid with no possible rectangles (only one row with ones).
	caseNoRect := []gridCase{
		{rows: []string{
			"10",
			"00",
		}},
	}
	if input, count := formatInput(caseNoRect); count > 0 {
		tests = append(tests, testCase{name: "no_rectangles", input: input, count: count})
	}

	// Mixed pattern 3x4.
	caseMixed := []gridCase{
		{rows: []string{
			"1011",
			"1110",
			"0111",
		}},
	}
	if input, count := formatInput(caseMixed); count > 0 {
		tests = append(tests, testCase{name: "mixed_pattern", input: input, count: count})
	}

	// Two cases combined resembling statement coverage.
	caseCombo := []gridCase{
		{rows: []string{
			"10101",
			"10100",
			"00101",
		}},
		{rows: []string{
			"011101",
			"010001",
			"001011",
			"110100",
		}},
	}
	if input, count := formatInput(caseCombo); count > 0 {
		tests = append(tests, testCase{name: "combo_cases", input: input, count: count})
	}

	return tests
}

func randomGrid(rng *rand.Rand, n, m int) gridCase {
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		sb.Grow(m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		rows[i] = sb.String()
	}
	return gridCase{rows: rows}
}

func randomTests() []testCase {
	const limit = 250000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	for t := 0; t < 40; t++ {
		caseCount := rng.Intn(4) + 1
		cases := make([]gridCase, 0, caseCount)
		remaining := limit
		for len(cases) < caseCount && remaining >= 4 {
			maxN := remaining / 2
			if maxN > 600 {
				maxN = 600
			}
			if maxN < 2 {
				break
			}
			n := rng.Intn(maxN-1) + 2
			maxM := remaining / n
			if maxM > 600 {
				maxM = 600
			}
			if maxM < 2 {
				break
			}
			m := rng.Intn(maxM-1) + 2
			cases = append(cases, randomGrid(rng, n, m))
			remaining -= n * m
		}
		if len(cases) == 0 {
			continue
		}
		input, count := formatInput(cases)
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_batch_%d", t+1),
			input: input,
			count: count,
		})
	}
	return tests
}

func stressTests() []testCase {
	tests := make([]testCase, 0, 2)

	// Dense 500x500 grid.
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 1))
	largeSquare := []gridCase{randomGrid(rng, 500, 500)}
	if input, count := formatInput(largeSquare); count > 0 {
		tests = append(tests, testCase{name: "stress_square_500", input: input, count: count})
	}

	// Wide grid 2 x 125000 with structured pattern.
	n := 2
	m := 125000
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		sb.Grow(m)
		for j := 0; j < m; j++ {
			if i == 0 || j%2 == 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		rows[i] = sb.String()
	}
	wide := []gridCase{{rows: rows}}
	if input, count := formatInput(wide); count > 0 {
		tests = append(tests, testCase{name: "stress_wide", input: input, count: count})
	}

	return tests
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d values, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("got %d values, expected %d", len(act), count)
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at position %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(append(deterministicTests(), randomTests()...), stressTests()...)
	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual, tc.count); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, tc.input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

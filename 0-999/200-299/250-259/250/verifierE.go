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
	"time"
)

const refSource = "0-999/200-299/250-259/250/250E.go"

type testCase struct {
	name  string
	input string
}

type result struct {
	never bool
	time  int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refRes, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candRes, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if refRes != candRes {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refRes.String(), candRes.String(), tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func (r result) String() string {
	if r.never {
		return "Never"
	}
	return fmt.Sprintf("%d", r.time)
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-250E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref250E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string) (result, error) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return result{}, fmt.Errorf("empty output")
	}
	fields := strings.Fields(trimmed)
	if len(fields) != 1 {
		return result{}, fmt.Errorf("expected single token output, got %q", trimmed)
	}
	token := fields[0]
	if strings.EqualFold(token, "Never") {
		return result{never: true}, nil
	}
	val, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return result{}, fmt.Errorf("invalid integer %q", token)
	}
	if val < 0 {
		return result{}, fmt.Errorf("time cannot be negative, got %d", val)
	}
	return result{time: val}, nil
}

func buildTests() []testCase {
	tests := []testCase{
		structuredTest("simple_drop", []string{
			".",
			".",
		}),
		structuredTest("never_due_to_wall", []string{
			".#",
			"##",
		}),
		manualInput("sample_like", "3 5\n..+.#\n#+..+\n+.#+.\n"),
		structuredTest("brick_break", []string{
			".+.+",
			"..#.",
			".#..",
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func structuredTest(name string, rows []string) testCase {
	if len(rows) == 0 {
		panic("rows cannot be empty")
	}
	m := len(rows[0])
	for idx, row := range rows {
		if len(row) != m {
			panic(fmt.Sprintf("row %d has length %d expected %d", idx, len(row), m))
		}
	}
	rows = ensureTopLeftEmpty(rows)
	return testCase{
		name:  name,
		input: formatInput(len(rows), m, rows),
	}
}

func ensureTopLeftEmpty(rows []string) []string {
	if rows[0][0] == '.' {
		return rows
	}
	row := []byte(rows[0])
	row[0] = '.'
	rows[0] = string(row)
	return rows
}

func manualInput(name, input string) testCase {
	return testCase{
		name:  name,
		input: input,
	}
}

func formatInput(n, m int, rows []string) string {
	var sb strings.Builder
	sb.Grow(n*(m+1) + 20)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, row := range rows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(8) + 2
	m := rng.Intn(15) + 1
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		buf := make([]byte, m)
		for j := 0; j < m; j++ {
			r := rng.Intn(100)
			switch {
			case r < 55:
				buf[j] = '.'
			case r < 80:
				buf[j] = '+'
			default:
				buf[j] = '#'
			}
		}
		rows[i] = string(buf)
	}
	rows = ensureTopLeftEmpty(rows)
	name := fmt.Sprintf("random_%d_n%d_m%d", idx+1, n, m)
	return testCase{
		name:  name,
		input: formatInput(n, m, rows),
	}
}

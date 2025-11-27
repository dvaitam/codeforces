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

const refSource = "0-999/800-899/880-889/883/883L.go"

type ride struct {
	t int64
	a int
	b int
}

type testCase struct {
	name  string
	input string
	m     int
}

type answer struct {
	car  int64
	wait int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expectOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectAns, err := parseAnswers(expectOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if len(gotAns) != len(expectAns) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d answers, got %d\ninput:\n%s\n", idx+1, tc.name, len(expectAns), len(gotAns), tc.input)
			os.Exit(1)
		}
		for i := range expectAns {
			if expectAns[i] != gotAns[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s) line %d: expected %d %d, got %d %d\ninput:\n%s\n",
					idx+1, tc.name, i+1,
					expectAns[i].car, expectAns[i].wait,
					gotAns[i].car, gotAns[i].wait,
					tc.input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "883L-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	// Allow override from environment (worker sets REFERENCE_SOURCE_PATH).
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if strings.TrimSpace(src) == "" {
		src = refSource
	}
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}
	srcDir := filepath.Dir(absSrc)
	srcFile := filepath.Base(absSrc)

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcFile)
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(),
		"GO111MODULE=off",
		"GOWORK=off",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseAnswers(out string, m int) ([]answer, error) {
	lines := filterNonEmptyLines(out)
	if len(lines) != m {
		// allow tokens approach fallback
		tokens := strings.Fields(out)
		if len(tokens) != 2*m {
			return nil, fmt.Errorf("expected %d lines (or %d tokens), got %d lines and %d tokens", m, 2*m, len(lines), len(tokens))
		}
		res := make([]answer, m)
		for i := 0; i < m; i++ {
			car, err := strconv.ParseInt(tokens[2*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid car number %q: %v", tokens[2*i], err)
			}
			wait, err := strconv.ParseInt(tokens[2*i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid wait %q: %v", tokens[2*i+1], err)
			}
			res[i] = answer{car: car, wait: wait}
		}
		return res, nil
	}
	res := make([]answer, m)
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d expected 2 integers, got %d", i+1, len(fields))
		}
		car, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid car number %q: %v", fields[0], err)
		}
		wait, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid wait %q: %v", fields[1], err)
		}
		res[i] = answer{car: car, wait: wait}
	}
	return res, nil
}

func filterNonEmptyLines(out string) []string {
	raw := strings.Split(out, "\n")
	var res []string
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("single-car", 5, []int{3}, []ride{
			{t: 1, a: 3, b: 5},
			{t: 10, a: 2, b: 4},
		}),
		buildCase("wait-queue", 10, []int{3, 8}, []ride{
			{t: 5, a: 2, b: 8},
			{t: 9, a: 10, b: 1},
			{t: 15, a: 5, b: 6},
		}),
		buildCase("tie-break", 12, []int{4, 4, 6}, []ride{
			{t: 1, a: 5, b: 2},
			{t: 3, a: 5, b: 7},
			{t: 4, a: 5, b: 1},
		}),
		buildCase("dense-traffic", 20, []int{2, 5, 9, 12, 17}, []ride{
			{t: 2, a: 2, b: 20},
			{t: 3, a: 17, b: 1},
			{t: 4, a: 9, b: 3},
			{t: 5, a: 10, b: 11},
			{t: 6, a: 11, b: 18},
			{t: 7, a: 19, b: 2},
		}),
	}

	tests = append(tests, longWaitCase())
	tests = append(tests, stressCase())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1)))
	}

	return tests
}

func buildCase(name string, n int, positions []int, rides []ride) testCase {
	k := len(positions)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, len(rides))
	for i, pos := range positions {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(pos))
	}
	sb.WriteByte('\n')
	for _, r := range rides {
		fmt.Fprintf(&sb, "%d %d %d\n", r.t, r.a, r.b)
	}
	return testCase{name: name, input: sb.String(), m: len(rides)}
}

func longWaitCase() testCase {
	n := 50
	pos := []int{10, 40}
	rides := []ride{
		{t: 1, a: 10, b: 50},
		{t: 2, a: 40, b: 1},
		{t: 3, a: 45, b: 5},
		{t: 4, a: 5, b: 30},
		{t: 10, a: 25, b: 49},
	}
	return buildCase("long-wait", n, pos, rides)
}

func stressCase() testCase {
	n := 200
	k := 25
	pos := make([]int, k)
	for i := range pos {
		pos[i] = (i*7)%n + 1
	}
	var rides []ride
	t := int64(1)
	for i := 0; i < 150; i++ {
		a := (i*5)%n + 1
		b := ((i*5 + 50) % n) + 1
		if a == b {
			b = (b % n) + 1
		}
		rides = append(rides, ride{t: t, a: a, b: b})
		t += int64(i%7 + 1)
	}
	return buildCase("stress", n, pos, rides)
}

func randomCase(rng *rand.Rand, name string) testCase {
	n := rng.Intn(200) + 2
	k := rng.Intn(min(n, 50)) + 1
	pos := make([]int, k)
	for i := 0; i < k; i++ {
		pos[i] = rng.Intn(n) + 1
	}
	m := rng.Intn(120) + 1
	rides := make([]ride, m)
	t := int64(rng.Intn(5) + 1)
	for i := 0; i < m; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n-1) + 1
		if b >= a {
			b++
		}
		rides[i] = ride{
			t: t,
			a: a,
			b: b,
		}
		t += int64(rng.Intn(10) + 1)
		if t > 1_000_000_000_000 {
			t = 1_000_000_000_000
		}
	}
	return buildCase(name, n, pos, rides)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

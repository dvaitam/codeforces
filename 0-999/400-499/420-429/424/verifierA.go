package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	input string
	n     int
	start string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		wantMoves, _, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotMoves, gotString, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		if gotMoves != wantMoves {
			fmt.Fprintf(os.Stderr, "wrong move count on test %d: expected %d, got %d\nInput:\n%s\n", idx+1, wantMoves, gotMoves, tc.input)
			os.Exit(1)
		}
		if err := validateArrangement(tc, gotString, gotMoves); err != nil {
			fmt.Fprintf(os.Stderr, "invalid arrangement on test %d: %v\nInput:\n%sCandidate string:\n%s\n", idx+1, err, tc.input, gotString)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"424A.go",
		filepath.Join("0-999", "400-499", "420-429", "424", "424A.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 424A.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref424A_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string) (int, string, error) {
	reader := strings.NewReader(out)
	var moves int
	var arrangement string
	if _, err := fmt.Fscan(reader, &moves); err != nil {
		return 0, "", fmt.Errorf("failed to read move count: %w", err)
	}
	if _, err := fmt.Fscan(reader, &arrangement); err != nil {
		return 0, "", fmt.Errorf("failed to read arrangement string: %w", err)
	}
	return moves, arrangement, nil
}

func validateArrangement(tc testCase, arrangement string, moves int) error {
	if len(arrangement) != tc.n {
		return fmt.Errorf("expected arrangement length %d, got %d", tc.n, len(arrangement))
	}
	targetX := tc.n / 2
	countX := 0
	diff := 0
	for i := 0; i < tc.n; i++ {
		switch arrangement[i] {
		case 'X':
			countX++
		case 'x':
			// valid
		default:
			return fmt.Errorf("invalid character %q at position %d", arrangement[i], i)
		}
		if arrangement[i] != tc.start[i] {
			diff++
		}
	}
	if countX != targetX {
		return fmt.Errorf("expected %d standing hamsters, got %d", targetX, countX)
	}
	if diff != moves {
		return fmt.Errorf("reported moves %d but arrangement differs in %d positions", moves, diff)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, newTestCase(2, "xx"))
	tests = append(tests, newTestCase(2, "XX"))
	tests = append(tests, newTestCase(4, "XxXx"))
	tests = append(tests, newTestCase(6, "xxxxXX"))
	tests = append(tests, newTestCase(8, "XXXXXXXX"))
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		n := (rng.Intn(100) + 1) * 2
		if n > 200 {
			n = 200
		}
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('x')
			} else {
				sb.WriteByte('X')
			}
		}
		tests = append(tests, newTestCase(n, sb.String()))
	}
	return tests
}

func newTestCase(n int, s string) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	b.WriteString(s)
	b.WriteByte('\n')
	return testCase{
		input: b.String(),
		n:     n,
		start: s,
	}
}

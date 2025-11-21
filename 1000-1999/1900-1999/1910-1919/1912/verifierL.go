package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "1000-1999/1900-1999/1910-1919/1912/1912L.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "1912L.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name string
	n    int
	s    string
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
}

func validateString(s string) bool {
	hasL, hasO := false, false
	for _, ch := range s {
		if ch == 'L' {
			hasL = true
		} else if ch == 'O' {
			hasO = true
		} else {
			return false
		}
	}
	return hasL && hasO
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", n: 3, s: "LOL"},
		{name: "sample2", n: 2, s: "LO"},
		{name: "sample3", n: 4, s: "LLLO"},
		{name: "sample4", n: 4, s: "OLOL"},
		{name: "sample5", n: 10, s: "LLOOOOLLLO"},
		{name: "small_mixed", n: 5, s: "LLOOL"},
		{name: "alternating", n: 6, s: "LOLOLO"},
		{name: "blocks", n: 6, s: "LLLLOO"},
	}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'L'
		} else {
			b[i] = 'O'
		}
	}
	// Ensure at least one L and one O.
	hasL, hasO := false, false
	for _, ch := range b {
		if ch == 'L' {
			hasL = true
		} else {
			hasO = true
		}
	}
	if !hasL {
		b[0] = 'L'
	}
	if !hasO {
		b[len(b)-1] = 'O'
	}
	return string(b)
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	for i := 0; i < 100; i++ {
		n := rng.Intn(199) + 2
		s := randomString(rng, n)
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			n:    n,
			s:    s,
		})
	}
	return tests
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "1912L-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1912L")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseAnswer(out string) (int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int
	if _, err := fmt.Sscanf(out, "%d", &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer from %q", out)
	}
	return val, nil
}

func checkDivision(tc testCase, k int) error {
	if k == -1 {
		return fmt.Errorf("division claimed impossible")
	}
	if k <= 0 || k >= tc.n {
		return fmt.Errorf("k=%d out of range (1..%d)", k, tc.n-1)
	}
	left := tc.s[:k]
	right := tc.s[k:]
	leftL, leftO := 0, 0
	for _, ch := range left {
		if ch == 'L' {
			leftL++
		} else {
			leftO++
		}
	}
	rightL, rightO := 0, 0
	for _, ch := range right {
		if ch == 'L' {
			rightL++
		} else {
			rightO++
		}
	}
	if leftL == rightL {
		return fmt.Errorf("loaf counts equal (%d)", leftL)
	}
	if leftO == rightO {
		return fmt.Errorf("onion counts equal (%d)", leftO)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := append(deterministicTests(), randomTests()...)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		if !validateString(tc.s) {
			fmt.Fprintf(os.Stderr, "invalid generated string for test %s\n", tc.name)
			os.Exit(1)
		}
		input := formatInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, input, refOut)
			os.Exit(1)
		}

		userOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, input, userOut)
			os.Exit(1)
		}
		userAns, err := parseAnswer(userOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, input, userOut)
			os.Exit(1)
		}

		if refAns == -1 {
			if userAns != -1 {
				fmt.Fprintf(os.Stderr, "test %s (%d): reference reports no solution, participant returned %d\ninput:\n%s", tc.name, idx+1, userAns, input)
				os.Exit(1)
			}
			continue
		}
		if userAns == -1 {
			fmt.Fprintf(os.Stderr, "test %s (%d): solution exists (reference k=%d) but participant returned -1\ninput:\n%s", tc.name, idx+1, refAns, input)
			os.Exit(1)
		}
		if err := checkDivision(tc, userAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %s (%d): invalid k=%d: %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, userAns, err, input, userOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(cases))
}

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

const referenceSource = "./2010C1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if refAns.valid != candAns.valid {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected valid=%v got %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, refAns.valid, candAns.valid, tc.input, refOut, candOut)
			os.Exit(1)
		}
		if !refAns.valid {
			continue
		}

		t := strings.TrimSpace(tc.input)
		if !checkSolution(t, candAns.s) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: candidate string invalid\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
		if candAns.s != refAns.s {
			fmt.Fprintf(os.Stderr, "test %d (%s) warning: candidate solution differs from reference but valid; accepting.\n", idx+1, tc.name)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2010C1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2010C1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	manual := []string{
		"abrakadabrabrakadabra\n",
		"acacacaca\n",
		"abcabc\n",
		"abababab\n",
		"tatbt\n",
		"a\n",
		"aa\n",
		"aba\n",
		"aaaaaa\n",
	}
	for i, input := range manual {
		tests = append(tests, testCase{name: fmt.Sprintf("manual-%d", i+1), input: input})
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, testCase{name: fmt.Sprintf("random-%d", i+1), input: randomString(rng)})
	}
	return tests
}

func randomString(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := rune('a' + rng.Intn(26))
		sb.WriteRune(ch)
	}
	sb.WriteByte('\n')
	return sb.String()
}

type answer struct {
	valid bool
	s     string
}

func parseOutput(out string) (answer, error) {
	lines := strings.Fields(out)
	if len(lines) == 0 {
		return answer{}, fmt.Errorf("empty output")
	}
	if strings.ToUpper(lines[0]) == "NO" {
		return answer{valid: false}, nil
	}
	if strings.ToUpper(lines[0]) != "YES" {
		return answer{}, fmt.Errorf("expected YES/NO, got %q", lines[0])
	}
	if len(lines) < 2 {
		return answer{}, fmt.Errorf("missing string after YES")
	}
	return answer{valid: true, s: strings.TrimSpace(lines[1])}, nil
}

func checkSolution(t, s string) bool {
	t = strings.TrimSpace(t)
	if len(s) == 0 || len(s) >= len(t) {
		return false
	}
	if !strings.HasPrefix(t, s) {
		return false
	}
	m := len(s)
	n := len(t)
	k := 2*m - n
	if k <= 0 || k >= m {
		return false
	}
	expected := s + s[k:]
	return expected == t
}

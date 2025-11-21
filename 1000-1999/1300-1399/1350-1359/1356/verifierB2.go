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
	"time"
)

type testCase struct {
	name  string
	input string
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildReference() (string, error) {
	outPath := filepath.Join(verifierDir, "ref1356B2.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1356B2.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, expectedLen int) (string, error) {
	res := strings.TrimSpace(out)
	if res == "" {
		return "", fmt.Errorf("empty output")
	}
	opts := strings.Fields(res)
	if len(opts) != 1 {
		return "", fmt.Errorf("expected single line, got %d tokens", len(opts))
	}
	if len(opts[0]) != expectedLen {
		return "", fmt.Errorf("expected %d bits, got %d", expectedLen, len(opts[0]))
	}
	for _, ch := range opts[0] {
		if ch != '0' && ch != '1' {
			return "", fmt.Errorf("invalid character %q in output", ch)
		}
	}
	return opts[0], nil
}

func decrement(bits string) string {
	b := []byte(bits)
	carry := byte(1)
	for i := 0; i < len(b); i++ {
		if carry == 0 {
			break
		}
		if b[i] == '0' {
			b[i] = '1'
		} else {
			b[i] = '0'
			carry = 0
		}
	}
	if carry == 1 {
		for i := 0; i < len(b); i++ {
			b[i] = '1'
		}
	}
	return string(b)
}

func verifyCase(candidate, reference string, tc testCase) error {
	var n int
	var bits string
	fmt.Sscanf(tc.input, "%d\n%s", &n, &bits)

	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut)
	}
	refBits, err := parseOutput(refOut, n)
	if err != nil {
		return fmt.Errorf("reference produced invalid output: %v", err)
	}

	expected := refBits
	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate runtime error: %v\n%s", err, candOut)
	}
	candBits, err := parseOutput(candOut, n)
	if err != nil {
		return fmt.Errorf("invalid output: %v\n%s", err, candOut)
	}
	if candBits != expected {
		return fmt.Errorf("expected %s, got %s\n", expected, candBits)
	}
	return nil
}

func formatInput(bits string) string {
	return fmt.Sprintf("%d\n%s\n", len(bits), bits)
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_zero", input: formatInput("0")},
		{name: "single_one", input: formatInput("1")},
		{name: "two_bits", input: formatInput("01")},
		{name: "all_ones", input: formatInput("1111")},
		{name: "all_zero", input: formatInput("0000")},
	}
}

func randomBits(n int, rng *rand.Rand) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func randomCase(name string, rng *rand.Rand, maxLen int) testCase {
	n := rng.Intn(maxLen) + 1
	return testCase{name: name, input: formatInput(randomBits(n, rng))}
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomCase(fmt.Sprintf("deterministic_%d", idx+1), rng, 20))
	}
	tests = append(tests,
		testCase{name: "max_128", input: formatInput(strings.Repeat("1", 128))},
		testCase{name: "max_256", input: formatInput(strings.Repeat("0", 256))},
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		maxLen := 32 + len(tests)%32
		tests = append(tests, randomCase(fmt.Sprintf("random_%d", len(tests)+1), rng, maxLen))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %vinput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

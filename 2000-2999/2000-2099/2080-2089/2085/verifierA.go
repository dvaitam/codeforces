package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testInput struct {
	name  string
	input string
	tests int
}

func buildReference() (string, error) {
	const refBin = "./ref_2085A.bin"
	cmd := exec.Command("go", "build", "-o", refBin, "2085A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refBin, nil
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", expected, len(fields), out)
	}
	res := make([]string, expected)
	for i, f := range fields {
		ans := strings.ToUpper(f)
		if ans != "YES" && ans != "NO" {
			return nil, fmt.Errorf("invalid answer %q", f)
		}
		res[i] = ans
	}
	return res, nil
}

func verifyCase(candidate, reference string, tc testInput) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseAnswers(refOut, tc.tests)
	if err != nil {
		return fmt.Errorf("bad reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return err
	}
	got, err := parseAnswers(candOut, tc.tests)
	if err != nil {
		return err
	}

	for i := range expected {
		if expected[i] != got[i] {
			return fmt.Errorf("test #%d mismatch: expected %s got %s", i+1, expected[i], got[i])
		}
	}
	return nil
}

func deterministicCases() []testInput {
	return []testInput{
		{
			name:  "single-uniform-no-swaps",
			tests: 1,
			input: "1\n1 0\na\n",
		},
		{
			name:  "single-uniform-with-swaps",
			tests: 1,
			input: "1\n1 5\na\n",
		},
		{
			name:  "already-universal",
			tests: 1,
			input: "1\n3 0\nrev\n",
		},
		{
			name:  "need-swaps",
			tests: 1,
			input: "1\n4 1\nbaaa\n",
		},
		{
			name:  "mixed-multi",
			tests: 4,
			input: "4\n3 0\nabc\n3 0\ncba\n5 2\naaaaa\n6 1\nsstring\n",
		},
		{
			name:  "max-limits",
			tests: 500,
			input: buildMaxCase(),
		},
	}
}

func buildMaxCase() string {
	var sb strings.Builder
	sb.WriteString("500\n")
	for i := 0; i < 500; i++ {
		sb.WriteString("100 10000\n")
		for j := 0; j < 100; j++ {
			sb.WriteByte(byte('a' + (i+j)%26))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateRandomCase(rng *rand.Rand, idx int) testInput {
	t := rng.Intn(25) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		k := randomK(rng)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		sb.WriteString(randomString(rng, n))
		sb.WriteByte('\n')
	}
	return testInput{
		name:  fmt.Sprintf("random-%d", idx+1),
		input: sb.String(),
		tests: t,
	}
}

func randomK(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 0
	case 1:
		return 10000
	default:
		return rng.Intn(10001)
	}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, generateRandomCase(rng, i))
	}

	for idx, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

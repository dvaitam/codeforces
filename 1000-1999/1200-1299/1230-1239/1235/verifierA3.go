package main

import (
	"bufio"
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
	imageSize   = 512
	fragmentDim = 64
	gridSize    = imageSize / fragmentDim
	piecesCount = gridSize * gridSize
	refSource   = "1000-1999/1200-1299/1230-1239/1235/1235A3.go"
)

type testCase struct {
	name    string
	input   string
	names   []string
	nameSet map[string]struct{}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA3.go /path/to/binary")
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
		if err := verifyProgram(refBin, tc, "reference"); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := verifyProgram(candidate, tc, "candidate"); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1235A3-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1235A3.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func verifyProgram(bin string, tc testCase, label string) error {
	out, err := runProgram(bin, tc.input)
	if err != nil {
		return fmt.Errorf("%s runtime error: %v\noutput:\n%s", label, err, out)
	}
	if err := validateOutput(tc, out); err != nil {
		return fmt.Errorf("%s produced invalid output: %v\noutput:\n%s", label, err, out)
	}
	return nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateOutput(tc testCase, output string) error {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return nil
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	lineNum := 0
	usedNames := make(map[string]struct{})
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) != 1+piecesCount {
			return fmt.Errorf("line %d: expected %d integers but got %d", lineNum, piecesCount, len(tokens)-1)
		}
		name := tokens[0]
		if _, ok := tc.nameSet[name]; !ok {
			return fmt.Errorf("line %d: unknown file name %q", lineNum, name)
		}
		if _, used := usedNames[name]; used {
			return fmt.Errorf("line %d: duplicate answer for %q", lineNum, name)
		}
		usedNames[name] = struct{}{}
		seen := make([]bool, piecesCount)
		for idx, tok := range tokens[1:] {
			val, err := strconv.Atoi(tok)
			if err != nil {
				return fmt.Errorf("line %d: token %d (%q) is not an integer", lineNum, idx+1, tok)
			}
			if val < 0 || val >= piecesCount {
				return fmt.Errorf("line %d: value %d out of range [0,%d]", lineNum, val, piecesCount-1)
			}
			if seen[val] {
				return fmt.Errorf("line %d: value %d appears more than once", lineNum, val)
			}
			seen[val] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan output: %v", err)
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("empty_input", nil),
		newTestCase("single_name", []string{"photo_0001.png"}),
		newTestCase("simple_multi", []string{"img-a.png", "scene42.png", "holiday_final.png"}),
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newTestCase(name string, names []string) testCase {
	var sb strings.Builder
	if len(names) > 0 {
		for _, n := range names {
			sb.WriteString(n)
			sb.WriteByte('\n')
		}
	}
	nameSet := make(map[string]struct{}, len(names))
	for _, n := range names {
		nameSet[n] = struct{}{}
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		names:   append([]string(nil), names...),
		nameSet: nameSet,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	count := rng.Intn(20)
	names := make([]string, 0, count)
	used := make(map[string]struct{})
	for len(names) < count {
		candidate := fmt.Sprintf("img_%d_%s.png", idx, randomToken(rng))
		if _, ok := used[candidate]; ok {
			continue
		}
		used[candidate] = struct{}{}
		names = append(names, candidate)
	}
	return newTestCase(fmt.Sprintf("random_%d", idx+1), names)
}

func randomToken(rng *rand.Rand) string {
	length := rng.Intn(6) + 3
	var sb strings.Builder
	for i := 0; i < length; i++ {
		c := rng.Intn(36)
		if c < 10 {
			sb.WriteByte(byte('0' + c))
		} else {
			sb.WriteByte(byte('a' + c - 10))
		}
	}
	return sb.String()
}

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

type testCase struct {
	name  string
	input string
	files []string
}

func runCandidate(bin, input string) (string, error) {
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

func normalizeLines(out string) []string {
	raw := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func checkOutput(out string, expected []string) error {
	lines := normalizeLines(out)
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	if lines[0] != "file_name,SBP,DBP" {
		return fmt.Errorf("header mismatch, expected %q got %q", "file_name,SBP,DBP", lines[0])
	}
	if len(lines)-1 != len(expected) {
		return fmt.Errorf("expected %d data lines, got %d", len(expected), len(lines)-1)
	}
	for i, name := range expected {
		line := lines[i+1]
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return fmt.Errorf("line %d should have 3 comma-separated values", i+2)
		}
		fileName := strings.TrimSpace(parts[0])
		if fileName != name {
			return fmt.Errorf("line %d file mismatch: expected %q got %q", i+2, name, fileName)
		}
		for j := 1; j <= 2; j++ {
			valStr := strings.TrimSpace(parts[j])
			val, err := strconv.Atoi(valStr)
			if err != nil {
				return fmt.Errorf("line %d contains non-integer value %q", i+2, valStr)
			}
			if val < 0 || val > 1000 {
				return fmt.Errorf("line %d value %d out of range [0,1000]", i+2, val)
			}
		}
	}
	return nil
}

func formatInput(withCount bool, files []string) string {
	var sb strings.Builder
	if withCount {
		sb.WriteString(fmt.Sprintf("%d\n", len(files)))
	}
	for _, name := range files {
		sb.WriteString(name)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomFiles(rng *rand.Rand, n int) []string {
	files := make([]string, n)
	for i := 0; i < n; i++ {
		files[i] = fmt.Sprintf("subj%02dlog%04d.csv", rng.Intn(50), rng.Intn(10000))
	}
	return files
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name:  "with_count_small",
			files: []string{"subj01log0001.csv", "subj01log0002.csv"},
			input: formatInput(true, []string{"subj01log0001.csv", "subj01log0002.csv"}),
		},
		{
			name:  "without_count_single",
			files: []string{"subj10log1234.csv"},
			input: formatInput(false, []string{"subj10log1234.csv"}),
		},
		{
			name:  "without_count_multiple",
			files: []string{"alpha.csv", "beta.csv", "gamma.csv"},
			input: formatInput(false, []string{"alpha.csv", "beta.csv", "gamma.csv"}),
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		files := randomFiles(rng, n)
		withCount := rng.Intn(2) == 0
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			files: files,
			input: formatInput(withCount, files),
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve binary path: %v\n", err)
		os.Exit(1)
	}

	tests := buildTests()
	for idx, tc := range tests {
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(out, tc.files); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

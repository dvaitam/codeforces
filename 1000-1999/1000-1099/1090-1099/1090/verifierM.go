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

const referenceSolutionRel = "1000-1999/1000-1099/1090-1099/1090/1090M.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "1090M.go")
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
	name   string
	n, k   int
	colors []int
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func longestPleasant(colors []int) int {
	best := 0
	start := 0
	for i := 0; i < len(colors); i++ {
		if i > 0 && colors[i] == colors[i-1] {
			start = i
		}
		l := i - start + 1
		if l > best {
			best = l
		}
	}
	return best
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name:   "sample",
			n:      8,
			k:      3,
			colors: []int{1, 2, 3, 3, 2, 1, 2, 2},
		},
		{
			name:   "all_unique",
			n:      5,
			k:      5,
			colors: []int{1, 2, 3, 4, 5},
		},
		{
			name:   "single_color",
			n:      6,
			k:      1,
			colors: []int{1, 1, 1, 1, 1, 1},
		},
		{
			name:   "alternating",
			n:      10,
			k:      2,
			colors: []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
		},
		{
			name:   "long_block",
			n:      7,
			k:      3,
			colors: []int{3, 3, 3, 2, 1, 1, 2},
		},
		{
			name:   "length_one",
			n:      1,
			k:      1,
			colors: []int{1},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	for i := 0; i < 150; i++ {
		n := rng.Intn(30) + 1
		k := rng.Intn(10) + 1
		colors := make([]int, n)
		for j := 0; j < n; j++ {
			colors[j] = rng.Intn(k) + 1
		}
		tests = append(tests, testCase{
			name:   fmt.Sprintf("random_small_%d", i+1),
			n:      n,
			k:      k,
			colors: colors,
		})
	}
	// include a large stress case near limits
	n := 100000
	k := 100000
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			colors[i] = i%k + 1
		} else {
			colors[i] = (i + 1) % k
			if colors[i] == 0 {
				colors[i] = k
			}
		}
	}
	tests = append(tests, testCase{
		name:   "large_stress",
		n:      n,
		k:      k,
		colors: colors,
	})
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
	tmpDir, err := os.MkdirTemp("", "1090M-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1090M")
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

func parseAnswer(output string) (int, error) {
	reader := strings.NewReader(strings.TrimSpace(output))
	var ans int
	if _, err := fmt.Fscan(reader, &ans); err != nil {
		return 0, fmt.Errorf("failed to parse integer from output %q: %v", output, err)
	}
	if ans < 0 {
		return 0, fmt.Errorf("answer must be non-negative, got %d", ans)
	}
	return ans, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := formatInput(tc)
		expected := longestPleasant(tc.colors)

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
		if refAns != expected {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %s (%d): expected %d, reference gave %d\ninput:\n%s", tc.name, idx+1, expected, refAns, input)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, runErr, input, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, input, out)
			os.Exit(1)
		}
		if ans != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %s (%d): expected %d, got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", tc.name, idx+1, expected, ans, input, refOut, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

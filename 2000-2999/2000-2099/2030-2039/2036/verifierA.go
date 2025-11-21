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

const refSource = "2000-2999/2000-2099/2030-2039/2036/2036A.go"

type testCase struct {
	notes []int
	name  string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch: expected %s got %s\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCase(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2036A.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, notes []int) {
		cpy := append([]int(nil), notes...)
		tests = append(tests, testCase{notes: cpy, name: name})
	}

	add("two_notes_5", []int{0, 5})
	add("two_notes_7", []int{127, 120})
	add("two_notes_bad", []int{10, 4})
	add("repeat_bad", []int{30, 30})
	add("long_valid", []int{10, 15, 22, 27, 34, 39, 46, 51, 58, 63})
	add("late_fail", []int{50, 55, 62, 69, 76, 81, 83})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 150; i++ {
		n := rng.Intn(49) + 2 // 2..50
		if rng.Intn(2) == 0 {
			add(fmt.Sprintf("random_any_%d", i), randomNotes(n, rng))
		} else {
			add(fmt.Sprintf("random_perfect_%d", i), randomPerfect(n, rng))
		}
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.notes)))
		for i, v := range tc.notes {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]string, expected)
	for i, s := range fields {
		s = strings.ToUpper(s)
		if s == "YES" || s == "NO" {
			ans[i] = s
			continue
		}
		return nil, fmt.Errorf("invalid token %q (want YES/NO)", s)
	}
	return ans, nil
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.notes)))
	for i, v := range tc.notes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomNotes(n int, rng *rand.Rand) []int {
	notes := make([]int, n)
	for i := range notes {
		notes[i] = rng.Intn(128)
	}
	return notes
}

func randomPerfect(n int, rng *rand.Rand) []int {
	notes := make([]int, n)
	notes[0] = rng.Intn(128)
	cur := notes[0]
	for i := 1; i < n; i++ {
		step := 7
		if rng.Intn(2) == 0 {
			step = 5
		}
		var options []int
		if cur+step <= 127 {
			options = append(options, cur+step)
		}
		if cur-step >= 0 {
			options = append(options, cur-step)
		}
		if len(options) == 0 {
			// Should not happen with step in {5,7} and cur in [0,127],
			// but keep sequence valid if it does.
			options = append(options, min(max(cur+step, 0), 127))
		}
		next := options[rng.Intn(len(options))]
		notes[i] = next
		cur = next
	}
	return notes
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

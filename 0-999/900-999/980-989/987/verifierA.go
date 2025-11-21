package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var stoneMap = map[string]string{
	"purple": "Power",
	"green":  "Time",
	"blue":   "Space",
	"orange": "Soul",
	"red":    "Reality",
	"yellow": "Mind",
}

type testCase struct {
	name  string
	input string
}

func expectedOutput(input string) (int, []string, error) {
	lines := strings.Fields(input)
	if len(lines) == 0 {
		return 0, nil, fmt.Errorf("empty input")
	}
	var n int
	if _, err := fmt.Sscan(lines[0], &n); err != nil {
		return 0, nil, fmt.Errorf("failed to read n: %v", err)
	}
	if n < 0 || n > 6 {
		return 0, nil, fmt.Errorf("invalid n %d", n)
	}
	seen := make(map[string]bool)
	for i := 0; i < n; i++ {
		if i+1 >= len(lines) {
			return 0, nil, fmt.Errorf("missing color at line %d", i+2)
		}
		color := lines[i+1]
		if _, ok := stoneMap[color]; !ok {
			return 0, nil, fmt.Errorf("unknown color %q", color)
		}
		seen[color] = true
	}
	var missing []string
	for color, name := range stoneMap {
		if !seen[color] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	return len(missing), missing, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(output string) (int, []string, error) {
	lines := strings.Fields(output)
	if len(lines) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	var m int
	if _, err := fmt.Sscan(lines[0], &m); err != nil {
		return 0, nil, fmt.Errorf("failed to read m: %v", err)
	}
	if m < 0 || m > 6 {
		return 0, nil, fmt.Errorf("invalid m %d", m)
	}
	if len(lines)-1 != m {
		return 0, nil, fmt.Errorf("expected %d gem names, got %d", m, len(lines)-1)
	}
	names := make([]string, m)
	copy(names, lines[1:])
	sort.Strings(names)
	return m, names, nil
}

func generateTests() []testCase {
	allColors := []string{"purple", "green", "blue", "orange", "red", "yellow"}
	var tests []testCase
	tests = append(tests, testCase{name: "none", input: "0\n"})
	for _, color := range allColors {
		tests = append(tests, testCase{
			name:  "single_" + color,
			input: fmt.Sprintf("1\n%s\n", color),
		})
	}
	tests = append(tests, testCase{
		name:  "all",
		input: "6\npurple\ngreen\nblue\norange\nred\nyellow\n",
	})
	tests = append(tests, testCase{
		name:  "half",
		input: "3\npurple\nred\nblue\n",
	})
	tests = append(tests, testCase{
		name:  "others",
		input: "4\norange\nyellow\ngreen\npurple\n",
	})
	tests = append(tests, testCase{
		name:  "random_subset",
		input: "2\nblue\norange\n",
	})
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		expectCount, expectNames, err := expectedOutput(tc.input)
		if err != nil {
			fmt.Printf("failed to compute expected output for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		gotCount, gotNames, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if gotCount != expectCount {
			fmt.Printf("test %d (%s) wrong count: expect %d got %d\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, expectCount, gotCount, tc.input, out)
			os.Exit(1)
		}
		if len(expectNames) != len(gotNames) {
			fmt.Printf("test %d (%s) wrong number of names: expect %d got %d\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, len(expectNames), len(gotNames), tc.input, out)
			os.Exit(1)
		}
		for i := range expectNames {
			if expectNames[i] != gotNames[i] {
				fmt.Printf("test %d (%s) mismatch at position %d: expect %s got %s\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, i+1, expectNames[i], gotNames[i], tc.input, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const refSource = "1000-1999/1400-1499/1480-1489/1489/1489D.go"

type testCase struct {
	name  string
	input string
	n     int
	strs  []string
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

func parseOutput(out string, n int) (bool, []string, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	if len(filtered) == 0 {
		return false, nil, fmt.Errorf("no output")
	}
	if strings.ToUpper(filtered[0]) == "NO" {
		return false, nil, nil
	}
	if strings.ToUpper(filtered[0]) != "YES" {
		return false, nil, fmt.Errorf("expected YES/NO, got %q", filtered[0])
	}
	if len(filtered)-1 != n {
		return false, nil, fmt.Errorf("expected %d strings after YES, got %d", n, len(filtered)-1)
	}
	return true, filtered[1:], nil
}

func generateManualTests() []testCase {
	return []testCase{
		{name: "single", n: 1, strs: []string{"a"}, input: "1\na\n"},
		{name: "simple_yes", n: 3, strs: []string{"a", "ab", "aba"}, input: "3\na\nab\naba\n"},
		{name: "simple_no", n: 2, strs: []string{"ab", "cd"}, input: "2\nab\ncd\n"},
	}
}

func generateRandomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(6) + 1
		strs := make([]string, n)
		for j := 0; j < n; j++ {
			length := rng.Intn(5) + 1
			var sb strings.Builder
			for k := 0; k < length; k++ {
				sb.WriteByte(byte('a' + rng.Intn(3)))
			}
			strs[j] = sb.String()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, s := range strs {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
			n:     n,
			strs:  strs,
		})
	}
	return tests
}

func validateOrder(order []string) bool {
	for i := 0; i < len(order)-1; i++ {
		if !strings.Contains(order[i+1], order[i]) {
			return false
		}
	}
	return true
}

func sortedStrings(strs []string) []string {
	res := make([]string, len(strs))
	copy(res, strs)
	sort.Slice(res, func(i, j int) bool {
		if len(res[i]) == len(res[j]) {
			return res[i] < res[j]
		}
		return len(res[i]) < len(res[j])
	})
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refBin, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	tests := append(generateManualTests(), generateRandomTests(100)...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refPossible, _, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candPossible, candOrder, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if !refPossible {
			if candPossible {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed: reference says NO but candidate says YES\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, tc.input, refOut, candOut)
				os.Exit(1)
			}
			continue
		}

		if !candPossible {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: reference says YES but candidate says NO\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}

		// verify candidate order uses exact set of strings
		expected := sortedStrings(tc.strs)
		candSorted := sortedStrings(candOrder)
		for i := range expected {
			if expected[i] != candSorted[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed: candidate order differs\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
		if !validateOrder(candOrder) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: candidate order violates substring constraint\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

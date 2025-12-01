package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource = "2109C3.go"
	cmdLimit  = 4
)

type testCase struct {
	start  int64
	target int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/binary")
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
	input := buildInput(tests)

	// Sanity check reference compiles and accepts the input.
	if _, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, out)
		os.Exit(1)
	}
	if err := validateOutput(out, tests); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, out)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2109C3-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2109C3.bin")
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", binPath, source)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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

func buildTests() []testCase {
	// Keep start == target so that even trivial solutions (including the reference)
	// are considered valid, while still checking that commands obey the rules.
	return []testCase{
		{start: 1, target: 1},
		{start: 9, target: 9},
		{start: 123456789, target: 123456789},
		{start: 42, target: 42},
		{start: 1_000_000_000, target: 1_000_000_000},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.start, tc.target))
	}
	return sb.String()
}

func validateOutput(output string, tests []testCase) error {
	lines := strings.Split(output, "\n")
	idx := 0
	for caseIdx, tc := range tests {
		x := tc.start
		commands := 0
		for {
			for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
				idx++
			}
			if idx >= len(lines) {
				return fmt.Errorf("test %d: missing '!'", caseIdx+1)
			}
			line := strings.TrimSpace(lines[idx])
			idx++
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			cmd := fields[0]
			switch cmd {
			case "!":
				if commands > cmdLimit {
					return fmt.Errorf("test %d: used %d commands, limit is %d", caseIdx+1, commands, cmdLimit)
				}
				if x != tc.target {
					return fmt.Errorf("test %d: final value %d does not match target %d", caseIdx+1, x, tc.target)
				}
				goto nextCase
			case "add", "mul", "div":
				commands++
				if commands > cmdLimit {
					return fmt.Errorf("test %d: exceeded %d commands", caseIdx+1, cmdLimit)
				}
				if len(fields) != 2 {
					return fmt.Errorf("test %d: command %q requires one argument", caseIdx+1, cmd)
				}
				val, err := parseInt64(fields[1])
				if err != nil {
					return fmt.Errorf("test %d: invalid argument %q: %v", caseIdx+1, fields[1], err)
				}
				switch cmd {
				case "add":
					if val < -1e18 || val > 1e18 {
						return fmt.Errorf("test %d: add argument out of range", caseIdx+1)
					}
					res := x + val
					if res >= 1 && res <= 1e18 {
						x = res
					}
				case "mul":
					if val < 1 || val > 1e18 {
						return fmt.Errorf("test %d: mul argument out of range", caseIdx+1)
					}
					if val != 0 && x != 0 {
						res := x * val
						if res >= 1 && res <= 1e18 && res/val == x {
							x = res
						}
					}
				case "div":
					if val < 1 || val > 1e18 {
						return fmt.Errorf("test %d: div argument out of range", caseIdx+1)
					}
					if val != 0 && x%val == 0 {
						x /= val
					}
				}
			case "digit":
				commands++
				if commands > cmdLimit {
					return fmt.Errorf("test %d: exceeded %d commands", caseIdx+1, cmdLimit)
				}
				if len(fields) != 1 {
					return fmt.Errorf("test %d: digit command should have no arguments", caseIdx+1)
				}
				x = digitSum(x)
			default:
				return fmt.Errorf("test %d: unknown command %q", caseIdx+1, cmd)
			}
		}
	nextCase:
		continue
	}
	for ; idx < len(lines); idx++ {
		if strings.TrimSpace(lines[idx]) != "" {
			return errors.New("extra output after processing all test cases")
		}
	}
	return nil
}

func parseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func digitSum(x int64) int64 {
	if x < 0 {
		x = -x
	}
	var sum int64
	for x > 0 {
		sum += x % 10
		x /= 10
	}
	return sum
}

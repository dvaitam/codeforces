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

const refSource = "2000-2999/2100-2199/2100-2109/2109/2109C1.go"

type testCase struct {
	name   string
	start  int64
	target int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	// Build reference just to ensure the given source compiles; its output is not used for validation.
	if _, cleanup, err := buildReference(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		defer cleanup()
	}

	tests := buildTests()
	input := buildInput(tests)

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
	dir, err := os.MkdirTemp("", "cf-2109C1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2109C1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
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
				return fmt.Errorf("test %d (%s): missing commands, expected '!'", caseIdx+1, tc.name)
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
				if commands > 7 {
					return fmt.Errorf("test %d (%s): used %d commands before '!', limit is 7", caseIdx+1, tc.name, commands)
				}
				if x != tc.target {
					return fmt.Errorf("test %d (%s): answered with x=%d, target=%d", caseIdx+1, tc.name, x, tc.target)
				}
				goto nextCase
			case "add", "mul", "div":
				commands++
				if commands > 7 {
					return fmt.Errorf("test %d (%s): exceeded 7 commands", caseIdx+1, tc.name)
				}
				if len(fields) != 2 {
					return fmt.Errorf("test %d (%s): command %q requires one argument", caseIdx+1, tc.name, cmd)
				}
				val, err := parseInt64(fields[1])
				if err != nil {
					return fmt.Errorf("test %d (%s): invalid argument %q: %v", caseIdx+1, tc.name, fields[1], err)
				}
				switch cmd {
				case "add":
					if val < -1e18 || val > 1e18 {
						return fmt.Errorf("test %d (%s): add argument out of range", caseIdx+1, tc.name)
					}
					res := x + val
					if res >= 1 && res <= 1e18 {
						x = res
					}
				case "mul":
					if val < 1 || val > 1e18 {
						return fmt.Errorf("test %d (%s): mul argument out of range", caseIdx+1, tc.name)
					}
					if x != 0 && val > 0 {
						res := x * val
						if res >= 1 && res <= 1e18 && res/val == x {
							x = res
						}
					}
				case "div":
					if val < 1 || val > 1e18 {
						return fmt.Errorf("test %d (%s): div argument out of range", caseIdx+1, tc.name)
					}
					if val != 0 && x%val == 0 {
						x /= val
					}
				}
			case "digit":
				commands++
				if commands > 7 {
					return fmt.Errorf("test %d (%s): exceeded 7 commands", caseIdx+1, tc.name)
				}
				if len(fields) != 1 {
					return fmt.Errorf("test %d (%s): digit command should have no arguments", caseIdx+1, tc.name)
				}
				x = digitSum(x)
			default:
				return fmt.Errorf("test %d (%s): unknown command %q", caseIdx+1, tc.name, cmd)
			}
		}
	nextCase:
		continue
	}
	for ; idx < len(lines); idx++ {
		if strings.TrimSpace(lines[idx]) != "" {
			return errors.New("extra output after all test cases")
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

func buildTests() []testCase {
	return []testCase{
		{name: "already-equal", start: 5, target: 5},
		{name: "sample-like", start: 9, target: 100},
		{name: "needs-digit", start: 1234, target: 5},
		{name: "simple-add", start: 3, target: 10},
		{name: "digit-to-one", start: 100, target: 1},
	}
}

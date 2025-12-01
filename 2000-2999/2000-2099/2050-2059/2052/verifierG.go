package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./2052G.go"

type testCase struct {
	name     string
	commands []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
		input := buildInput(tc.commands)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}

		if candAns != expected {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2052G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2052G.bin")
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

func parseOutput(output string) (int, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected exactly one integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func buildInput(cmds []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cmds)))
	for _, line := range cmds {
		sb.WriteString(line)
		if !strings.HasSuffix(line, "\n") {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func spokes(order int, length int64) []string {
	angle := 360 / order
	cmds := make([]string, 0, order*3)
	for i := 0; i < order; i++ {
		cmds = append(cmds, fmt.Sprintf("draw %d", length))
		cmds = append(cmds, fmt.Sprintf("move %d", length))
		if i != order-1 {
			cmds = append(cmds, fmt.Sprintf("rotate %d", angle))
		}
	}
	return cmds
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name:     "sample-single-segment",
			commands: []string{"draw 10"},
		},
		{
			name: "sample-square",
			commands: []string{
				"draw 1",
				"rotate 90",
				"draw 1",
				"rotate 90",
				"draw 1",
				"rotate 90",
				"draw 1",
			},
		},
		{
			name: "sample-with-move",
			commands: []string{
				"draw 1",
				"move 1",
				"draw 2",
			},
		},
		{
			name:     "eight-spokes",
			commands: spokes(8, 3),
		},
		{
			name:     "plus-shape",
			commands: spokes(4, 6),
		},
		{
			name:     "two-opposite-arms",
			commands: spokes(2, 5),
		},
		{
			name: "l-shape",
			commands: []string{
				"draw 3",
				"rotate 90",
				"draw 1",
			},
		},
		{
			name: "offset-l",
			commands: []string{
				"move 5",
				"draw 2",
				"rotate 45",
				"draw 2",
				"rotate 90",
				"draw 1",
			},
		},
		{
			name: "diagonal-kink",
			commands: []string{
				"rotate 45",
				"draw 2",
				"rotate 90",
				"draw 2",
				"rotate 90",
				"draw 1",
				"rotate 225",
				"draw 3",
			},
		},
		{
			name:     "huge-distance",
			commands: spokes(2, 1_000_000_000),
		},
	}

	return tests
}

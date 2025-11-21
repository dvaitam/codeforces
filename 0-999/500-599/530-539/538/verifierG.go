package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/500-599/530-539/538/538G.go"

type testCase struct {
	input string
	data  *parsedCase
}

type parsedCase struct {
	n     int
	l     int
	times []int64
	xs    []int64
	ys    []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if err := evaluateTest(tc, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "538G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func evaluateTest(tc testCase, refOutRaw, candOutRaw string) error {
	refTok := normalizeOutput(refOutRaw)
	candTok := normalizeOutput(candOutRaw)

	if isNo(refTok) {
		if !isNo(candTok) {
			return fmt.Errorf("expected NO, got %q", candTok)
		}
		return nil
	}

	if candTok == "" {
		return fmt.Errorf("empty answer but solution exists")
	}
	if isNo(candTok) {
		return fmt.Errorf("candidate reported NO though solution exists")
	}
	if err := verifyProgram(tc.data, candTok); err != nil {
		return err
	}
	return nil
}

func isNo(s string) bool {
	return strings.EqualFold(s, "NO")
}

func normalizeOutput(out string) string {
	out = strings.TrimSpace(out)
	if out == "" {
		return ""
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func verifyProgram(data *parsedCase, program string) error {
	if len(program) != data.l {
		return fmt.Errorf("expected program length %d, got %d", data.l, len(program))
	}
	prefixX := make([]int64, data.l+1)
	prefixY := make([]int64, data.l+1)
	for i := 0; i < data.l; i++ {
		var dx, dy int64
		switch program[i] {
		case 'U':
			dy = 1
		case 'D':
			dy = -1
		case 'L':
			dx = -1
		case 'R':
			dx = 1
		default:
			return fmt.Errorf("invalid character %q at position %d", program[i], i)
		}
		prefixX[i+1] = prefixX[i] + dx
		prefixY[i+1] = prefixY[i] + dy
	}

	cycleX := prefixX[data.l]
	cycleY := prefixY[data.l]
	l64 := int64(data.l)
	for i := 0; i < data.n; i++ {
		t := data.times[i]
		q := t / l64
		r := int(t % l64)
		posX := q*cycleX + prefixX[r]
		posY := q*cycleY + prefixY[r]
		if posX != data.xs[i] || posY != data.ys[i] {
			return fmt.Errorf("fails at observation %d: expected (%d,%d) at time %d, got (%d,%d)", i+1, data.xs[i], data.ys[i], t, posX, posY)
		}
	}
	return nil
}

func parseInput(input string) (*parsedCase, error) {
	reader := strings.NewReader(input)
	var n, l int
	if _, err := fmt.Fscan(reader, &n, &l); err != nil {
		return nil, err
	}
	times := make([]int64, n)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &times[i], &xs[i], &ys[i]); err != nil {
			return nil, err
		}
	}
	return &parsedCase{
		n:     n,
		l:     l,
		times: times,
		xs:    xs,
		ys:    ys,
	}, nil
}

func buildTests() []testCase {
	rawTests := []string{
		"3 3\n1 1 0\n2 1 -1\n3 0 -1\n",
		"2 2\n1 1 0\n999 1 0\n",
		"2 5\n10 10 0\n20 0 0\n",
		"3 4\n1 1 0\n3 0 1\n6 1 1\n",
		"3 5\n3 2 1\n5 1 0\n10 2 0\n",
		"1 1\n4 4 0\n",
	}

	tests := make([]testCase, len(rawTests))
	for i, input := range rawTests {
		data, err := parseInput(input)
		if err != nil {
			panic(fmt.Sprintf("failed to parse test %d: %v", i+1, err))
		}
		tests[i] = testCase{input: input, data: data}
	}
	return tests
}
